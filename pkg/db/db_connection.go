package db

import (
	"fmt"
	config "scs-user/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn(cfg)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return db, nil
}

func dsn(c *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s  sslmode=disable",
		c.Database.DbHost, c.Database.DbPort, c.Database.DbUser, c.Database.DbName, c.Database.DbPassword)
}
