package handler

import (
	"errors"
	"log/slog" // Добавляем slog
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nakle1ka/piggy/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserHandler struct {
	service service.UserService
}

func (h *UserHandler) Register(c *gin.Context) {
	var req CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("register: invalid json", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	accessToken, refreshToken, err := h.service.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			slog.Info("register: user already exists", slog.String("email", req.Email))
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		} else {
			slog.Error("register: server error", slog.Any("err", err), slog.String("email", req.Email))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	slog.Info("user registered", slog.String("email", req.Email))
	c.SetCookie("refresh_token", refreshToken, int(h.service.GetRefreshExpires().Seconds()), "/", "", true, true)
	c.JSON(http.StatusCreated, TokenDTO{AccessToken: accessToken})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("login: invalid json", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, service.ErrInvalidCredentials) {
			slog.Info("login: failed attempt", slog.String("email", req.Email))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		} else {
			slog.Error("login: server error", slog.Any("err", err), slog.String("email", req.Email))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	slog.Info("user logged in", slog.String("email", req.Email))
	c.SetCookie("refresh_token", refreshToken, int(h.service.GetRefreshExpires().Seconds()), "/", "", true, true)
	c.JSON(http.StatusOK, TokenDTO{AccessToken: accessToken})
}

func (h *UserHandler) Logout(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = h.service.LogoutUser(rt)
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) ||
			errors.Is(err, redis.Nil) ||
			errors.Is(err, jwt.ErrTokenExpired) {
			slog.Warn("logout: invalid or expired token", slog.Any("err", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		} else {
			slog.Error("logout: redis error", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error: could not delete session"})
		}
		return
	}

	slog.Info("user logged out")
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.Status(http.StatusOK)
}

func (h *UserHandler) Refresh(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.service.GetNewAccessToken(rt)
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) ||
			errors.Is(err, redis.Nil) ||
			errors.Is(err, jwt.ErrTokenExpired) {
			slog.Warn("refresh: token validation failed", slog.Any("err", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		} else {
			slog.Error("refresh: server error", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error: could not refresh token"})
		}
		return
	}

	slog.Info("token refreshed")
	c.SetCookie("refresh_token", refreshToken, int(h.service.GetRefreshExpires().Seconds()), "/", "", true, true)
	c.JSON(http.StatusOK, TokenDTO{AccessToken: accessToken})
}

func NewUserHandler(srv service.UserService) *UserHandler {
	return &UserHandler{
		service: srv,
	}
}
