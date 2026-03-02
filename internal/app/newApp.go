package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/piggy/internal/config"
	"github.com/nakle1ka/piggy/internal/handler"
	"github.com/nakle1ka/piggy/internal/middleware"
	"github.com/nakle1ka/piggy/internal/pkg/auth"
	"github.com/nakle1ka/piggy/internal/pkg/hash"
	"github.com/nakle1ka/piggy/internal/repository"
	"github.com/nakle1ka/piggy/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	cfg   *config.Config
	db    *gorm.DB
	cache *redis.Client
}

func (a *App) Run() error {
	tokenHasher := hash.NewTokenHasher()
	passwordHasher := hash.NewPasswordHasher()
	tokenManager := auth.NewTokenManager(a.cfg.JWT.SecretKey)

	userRepo := repository.NewUserRepository(a.db)
	piggyRepo := repository.NewPiggyRepository(a.db)
	cacheRepo := repository.NewCacheRepository(a.cache)

	userSrv := service.NewUserService(
		userRepo,
		cacheRepo,
		tokenHasher,
		passwordHasher,
		tokenManager,

		service.WithAccessExpires(a.cfg.JWT.AccessTokenLifeTime),
		service.WithRefreshExpires(a.cfg.JWT.RefreshTokenLifeTime),
	)
	piggySrv := service.NewPiggyService(piggyRepo)

	userHnd := handler.NewUserHandler(userSrv)
	piggyHnd := handler.NewPiggyHandler(piggySrv)

	router := gin.Default()

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHnd.Register)
		auth.POST("/login", userHnd.Login)
		auth.POST("/logout", userHnd.Logout)
		auth.POST("/refresh", userHnd.Refresh)
	}

	protected := v1.Group("/")
	protected.Use(middleware.JWTAuth(tokenManager))

	piggy := protected.Group("/piggies")
	{
		piggy.POST("/", piggyHnd.CreatePiggy)
		piggy.GET("/", piggyHnd.GetMyPiggies)
		piggy.GET("/:id", piggyHnd.GetPiggyById)
		piggy.POST("/:id/deposit", piggyHnd.Deposit)
		piggy.POST("/:id/withdrawal", piggyHnd.Withdrawal)
		piggy.PATCH("/:id", piggyHnd.UpdatePiggy)
		piggy.DELETE("/:id", piggyHnd.DeletePiggy)
	}

	protected.GET("/users/me", userHnd.GetMe)

	addr := fmt.Sprintf(":%v", a.cfg.App.Port)
	return router.Run(addr)
}

func NewApp(cfg *config.Config, db *gorm.DB, cache *redis.Client) *App {
	return &App{
		cfg:   cfg,
		db:    db,
		cache: cache,
	}
}
