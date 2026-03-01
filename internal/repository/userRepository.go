package repository

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/piggy/internal/model"
	"gorm.io/gorm"
)

type User = model.User
type UpdateUserDto struct {
	Username string
}

type UserRepository interface {
	CreateUser(user *User) error
	GetUserById(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(id uuid.UUID, user UpdateUserDto) error
	DeleteUserById(id uuid.UUID) error

	GetTransaction() *gorm.DB
	CreateUserWithTransaction(tx *gorm.DB, user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *userRepository) CreateUserWithTransaction(tx *gorm.DB, user *User) error {
	return tx.Create(user).Error
}

func (r *userRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) DeleteUserById(id uuid.UUID) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

func (r *userRepository) GetUserById(id uuid.UUID) (*User, error) {
	var user User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.First(&user, "email = ?", email).Error
	return &user, err
}

func (r *userRepository) UpdateUser(id uuid.UUID, user UpdateUserDto) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(user).Error
}

func NewUserRepository(pg *gorm.DB) UserRepository {
	return &userRepository{
		db: pg,
	}
}
