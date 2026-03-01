package service

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/piggy/internal/model"
	"github.com/nakle1ka/piggy/internal/repository"
)

type Piggy = model.Piggy

type PiggyService interface {
	CreatePiggy(piggy *Piggy) error
	UpdatePiggy(id, userId uuid.UUID, piggy Piggy) error
	DeletePiggy(id, userId uuid.UUID) error

	GetPiggyById(id, userId uuid.UUID) (*Piggy, error)
	GetPiggiesList(userId uuid.UUID, filters repository.PiggyFilters) ([]Piggy, error)

	Withdrawal(id, userId uuid.UUID, amount int64) error
	Deposit(id, userId uuid.UUID, amount int64) error
}

type piggyService struct {
	piggyRepo repository.PiggyRepository
}

func (s *piggyService) CreatePiggy(piggy *Piggy) error {
	return s.piggyRepo.CreatePiggy(piggy)
}

func (s *piggyService) DeletePiggy(id, userId uuid.UUID) error {
	return s.piggyRepo.DeletePiggy(id, userId)
}

func (s *piggyService) Deposit(id, userId uuid.UUID, amount int64) error {
	piggy, err := s.piggyRepo.GetPiggyById(id, userId)
	if err != nil {
		return err
	}
	if piggy.Accumulated+amount > piggy.Amount {
		return ErrInvalidAmount
	}
	return s.piggyRepo.AddPiggyAccumulated(id, userId, amount)
}

func (s *piggyService) Withdrawal(id, userId uuid.UUID, amount int64) error {
	piggy, err := s.piggyRepo.GetPiggyById(id, userId)
	if err != nil {
		return err
	}
	if piggy.Accumulated-amount < 0 {
		return ErrInvalidAmount
	}
	return s.piggyRepo.AddPiggyAccumulated(id, userId, -amount)
}

func (s *piggyService) GetPiggyById(id, useid uuid.UUID) (*Piggy, error) {
	return s.piggyRepo.GetPiggyById(id, useid)
}

func (s *piggyService) GetPiggiesList(userId uuid.UUID, filters repository.PiggyFilters) ([]Piggy, error) {
	return s.piggyRepo.GetPiggiesList(userId, filters)
}

func (s *piggyService) UpdatePiggy(id, userId uuid.UUID, piggy Piggy) error {
	return s.piggyRepo.UpdatePiggy(id, userId, piggy)
}

func NewPiggyService(r repository.PiggyRepository) PiggyService {
	return &piggyService{
		piggyRepo: r,
	}
}
