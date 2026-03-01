package repository

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/piggy/internal/model"
	"gorm.io/gorm"
)

type Piggy = model.Piggy

type PiggyFilters struct {
	Title *string
}

type PiggyRepository interface {
	CreatePiggy(piggy *Piggy) error
	DeletePiggy(id, userId uuid.UUID) error

	GetPiggyById(id, userId uuid.UUID) (*Piggy, error)
	GetPiggiesList(userId uuid.UUID, filters PiggyFilters) ([]Piggy, error)

	AddPiggyAccumulated(id, userId uuid.UUID, amount int64) error
	UpdatePiggy(id, userId uuid.UUID, piggy Piggy) error
}

type piggyRepository struct {
	db *gorm.DB
}

func (r *piggyRepository) GetPiggiesList(userId uuid.UUID, filters PiggyFilters) ([]Piggy, error) {
	var piggies []Piggy

	query := r.db.Where(
		"user_id = ?",
		userId,
	)

	if filters.Title != nil {
		query = query.Where("title ILIKE ?", "%"+*filters.Title+"%")
	}

	err := query.Find(&piggies).Error
	if err != nil {
		return nil, err
	}

	return piggies, nil
}

func (r *piggyRepository) GetPiggyById(id, userId uuid.UUID) (*Piggy, error) {
	var piggy Piggy

	err := r.db.
		First(&piggy, "id = ? AND user_id = ?", id, userId).Error
	if err != nil {
		return nil, err
	}

	return &piggy, nil
}

func (r *piggyRepository) CreatePiggy(piggy *Piggy) error {
	return r.db.Create(piggy).Error
}

func (r *piggyRepository) DeletePiggy(id, userId uuid.UUID) error {
	return r.db.
		Delete(Piggy{}, "id = ? AND user_id = ?", id, userId).Error
}

func (r *piggyRepository) UpdatePiggy(id, userId uuid.UUID, piggy Piggy) error {
	err := r.db.
		Model(&Piggy{}).
		Where("id = ? AND user_id = ?", id, userId).
		Updates(piggy).Error
	return err
}

func (r *piggyRepository) AddPiggyAccumulated(id, userId uuid.UUID, amount int64) error {
	return r.db.
		Model(&Piggy{}).
		Where("id = ? AND user_id = ?", id, userId).
		Update("accumulated", gorm.Expr("accumulated + ?", amount)).Error
}

func NewPiggyRepository(db *gorm.DB) PiggyRepository {
	return &piggyRepository{
		db: db,
	}
}
