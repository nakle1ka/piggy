package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/piggy/internal/model"
	"github.com/nakle1ka/piggy/internal/pkg/auth"
	"github.com/nakle1ka/piggy/internal/pkg/hash"
	"github.com/nakle1ka/piggy/internal/repository"
)

type UserService interface {
	RegisterUser(username, email, password string) (string, string, error)
	LoginUser(email, password string) (string, string, error)
	LogoutUser(refreshToken string) error
	GetNewAccessToken(refreshToken string) (string, string, error)
	GetUserById(id uuid.UUID) (*model.User, error)

	GetRefreshExpires() time.Duration
	GetAccessExpires() time.Duration
}

type userService struct {
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository

	passwordHasher hash.Hasher
	tokenHasher    hash.Hasher

	tokenManager auth.TokenManager

	refreshExp time.Duration
	accessExp  time.Duration
}

func (s *userService) GetUserById(id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetUserById(id)
}

func (s *userService) RegisterUser(username, email, password string) (string, string, error) {
	hashedPassword, err := s.passwordHasher.Hash([]byte(password))
	if err != nil {
		return "", "", err
	}

	tx := s.userRepo.GetTransaction()
	defer tx.Rollback()

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.CreateUserWithTransaction(tx, user); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := s.issueSession(user.Id)
	if err != nil {
		return "", "", ErrCreateToken
	}

	tx.Commit()
	return accessToken, refreshToken, nil
}

func (s *userService) LoginUser(email string, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	if !s.passwordHasher.Verify([]byte(password), []byte(user.PasswordHash)) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.issueSession(user.Id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) LogoutUser(refreshToken string) error {
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return err
	}
	userId, tokenId := claims.Subject, claims.TokenId

	key := fmt.Sprintf("session:%v:%v", userId, tokenId)
	return s.cacheRepo.Delete(key)
}

func (s *userService) GetNewAccessToken(refreshToken string) (string, string, error) {
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	userIdString, tokenId := claims.Subject, claims.TokenId
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return "", "", err
	}

	key := fmt.Sprintf("session:%v:%v", userId, tokenId)
	hash, err := s.cacheRepo.Get(key)
	if err != nil {
		return "", "", err
	}

	if !s.tokenHasher.Verify([]byte(refreshToken), []byte(hash)) {
		return "", "", ErrInvalidToken
	}

	accessToken, refreshToken, err := s.issueSession(userId)
	if err != nil {
		return "", "", err
	}

	err = s.cacheRepo.Delete(key)
	if err != nil {
		slog.Error("failed to delete session",
			slog.String("err", err.Error()),
			slog.String("key", key),
		)
	}

	return accessToken, refreshToken, nil
}

func (s *userService) issueSession(userId uuid.UUID) (string, string, error) {
	accessToken, _, err := s.tokenManager.GenerateToken(userId.String(), s.accessExp)
	if err != nil {
		return "", "", err
	}

	refreshToken, refreshTokenId, err := s.tokenManager.GenerateToken(userId.String(), s.refreshExp)
	if err != nil {
		return "", "", err
	}

	tokenHash, err := s.tokenHasher.Hash([]byte(refreshToken))
	if err != nil {
		return "", "", err
	}

	key := fmt.Sprintf("session:%v:%v", userId, refreshTokenId)
	if err := s.cacheRepo.Set(key, string(tokenHash), s.refreshExp); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) GetRefreshExpires() time.Duration {
	return s.refreshExp
}

func (s *userService) GetAccessExpires() time.Duration {
	return s.accessExp
}

type opt = func(s *userService)

func WithAccessExpires(exp int) opt {
	return func(s *userService) {
		s.accessExp = time.Duration(exp) * time.Second
	}
}

func WithRefreshExpires(exp int) opt {
	return func(s *userService) {
		s.refreshExp = time.Duration(exp) * time.Second
	}
}

func NewUserService(
	ur repository.UserRepository,
	cr repository.CacheRepository,
	th hash.Hasher,
	ph hash.Hasher,
	tm auth.TokenManager,
	opts ...opt,
) UserService {
	srv := &userService{
		userRepo:       ur,
		cacheRepo:      cr,
		tokenHasher:    th,
		passwordHasher: ph,
		tokenManager:   tm,

		refreshExp: time.Hour * 24 * 7,
		accessExp:  time.Minute * 15,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}
