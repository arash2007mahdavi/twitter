package database

import (
	"fmt"
	"time"
	"twitter/src/configs"
	"twitter/src/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logger.NewLogger()
var DBClient *gorm.DB

func InitDB(cfg *configs.Config) error {
	var err error

	dns := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v TimeZone=Asia/Tehran",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Dbname, cfg.Postgres.Sslmode,
	)

	DBClient, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}

	database, err := DBClient.DB()
	if err != nil {
		return err
	}

	err = database.Ping()
	if err != nil {
		return err
	}

	database.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifeTime * time.Minute)
	database.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	database.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)

	log.Info(logger.Postgres, logger.Start, "started successfuly", nil)
	return nil
}

func GetDB() *gorm.DB {
	return DBClient
}

func CloseDB() {
	database, _ := DBClient.DB()
	database.Close()
}