package storage

import (
	"fmt"
	"todo-service/src/models"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	zapgorm "moul.io/zapgorm2"
)

func InitPostgres(logger *zap.Logger) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.pass"),
		viper.GetString("db.name"),
		viper.GetString("db.port"),
	)
	log := zapgorm.New(logger)
	log.SetAsDefault()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	logger.Info("Migrating models")
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		panic(err)
	}

	return db
}
