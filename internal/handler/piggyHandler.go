package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/piggy/internal/model"
	"github.com/nakle1ka/piggy/internal/repository"
	"github.com/nakle1ka/piggy/internal/service"
	"gorm.io/gorm"
)

type PiggyHandler struct {
	service service.PiggyService
}

func (h *PiggyHandler) GetMyPiggies(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	log := slog.With("user_id", userId)

	var filters repository.PiggyFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		log.Warn("failed to bind piggy filters", slog.Any("err", err))
		c.Status(http.StatusBadRequest)
		return
	}

	piggies, err := h.service.GetPiggiesList(userId, filters)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			log.Error("failed to get piggies list", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not get piggies"})
		}
		return
	}

	c.JSON(http.StatusOK, piggies)
}

func (h *PiggyHandler) GetPiggyById(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	piggyId, ok := getPiggyId(c)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid piggy id",
		})
		return
	}
	log := slog.With("user_id", userId, "piggy_id", piggyId)

	piggy, err := h.service.GetPiggyById(piggyId, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			log.Error("failed to get piggy", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not get piggy"})
		}
		return
	}

	c.JSON(http.StatusOK, piggy)
}

func (h *PiggyHandler) CreatePiggy(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	log := slog.With("user_id", userId)

	var req CreatePiggyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	piggy := model.Piggy{
		UserId: userId,
		Title:  req.Title,
		Amount: req.Amount,
	}

	if err := h.service.CreatePiggy(&piggy); err != nil {
		log.Error("failed to create piggy", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create piggy"})
		return
	}

	log.Info("piggy created", slog.Any("piggy_id", piggy.Id))
	c.JSON(http.StatusOK, piggy)
}

func (h *PiggyHandler) UpdatePiggy(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	piggyId, ok := getPiggyId(c)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid piggy id",
		})
		return
	}
	log := slog.With("user_id", userId, "piggy_id", piggyId)

	var req UpdatePiggyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	piggy := model.Piggy{Title: req.Title}
	if err := h.service.UpdatePiggy(piggyId, userId, piggy); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			log.Error("failed to update piggy", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not update piggy"})
		}
		return
	}

	log.Info("piggy updated")
	c.Status(http.StatusOK)
}

func (h *PiggyHandler) DeletePiggy(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	piggyId, ok := getPiggyId(c)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid piggy id",
		})
		return
	}
	log := slog.With("user_id", userId, "piggy_id", piggyId)

	if err := h.service.DeletePiggy(piggyId, userId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			log.Error("failed to delete piggy", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not delete piggy"})
		}
		return
	}

	log.Info("piggy deleted")
	c.Status(http.StatusOK)
}

func (h *PiggyHandler) Deposit(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	piggyId, ok := getPiggyId(c)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid piggy id",
		})
		return
	}
	log := slog.With("user_id", userId, "piggy_id", piggyId)

	var req AmountPiggyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err := h.service.Deposit(piggyId, userId, req.Amount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
		} else if errors.Is(err, service.ErrInvalidAmount) {
			log.Warn("invalid deposit amount", slog.Int64("amount", req.Amount))
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid amount"})
		} else {
			log.Error("failed to deposit", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not deposit"})
		}
		return
	}

	log.Info("deposit successful", slog.Int64("amount", req.Amount))
	c.Status(http.StatusOK)
}

func (h *PiggyHandler) Withdrawal(c *gin.Context) {
	userId, ok := getUserId(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	piggyId, ok := getPiggyId(c)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid piggy id",
		})
		return
	}
	log := slog.With("user_id", userId, "piggy_id", piggyId)

	var req AmountPiggyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err := h.service.Withdrawal(piggyId, userId, req.Amount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusNotFound)
		} else if errors.Is(err, service.ErrInvalidAmount) {
			log.Warn("invalid withdrawal amount", slog.Int64("amount", req.Amount))
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid amount"})
		} else {
			log.Error("failed to withdraw", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not withdraw"})
		}
		return
	}

	log.Info("withdrawal successful", slog.Int64("amount", req.Amount))
	c.Status(http.StatusOK)
}

func NewPiggyHandler(srv service.PiggyService) *PiggyHandler {
	return &PiggyHandler{
		service: srv,
	}
}
