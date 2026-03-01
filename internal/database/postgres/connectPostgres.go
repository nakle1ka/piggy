package postgres

import (
	"fmt"

	"github.com/nakle1ka/piggy/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(opts ...opt) (*gorm.DB, error) {
	postgesCfg := defaultPostgresConfig()

	for _, opt := range opts {
		opt(&postgesCfg)
	}

	if postgesCfg.User == "" || postgesCfg.Password == "" || postgesCfg.DBname == "" {
		return nil, fmt.Errorf("postgres: invalid configuration")
	}

	var sslmode string
	if postgesCfg.Sslmode {
		sslmode = "prefer"
	} else {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		postgesCfg.Host,
		postgesCfg.User,
		postgesCfg.Password,
		postgesCfg.DBname,
		postgesCfg.Port,
		sslmode,
	)

	postges, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	migrationErr := postges.AutoMigrate(&model.User{}, &model.Piggy{})

	return postges, migrationErr
}
