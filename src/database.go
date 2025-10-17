package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDatabase(databaseURL string) error {
	var err error

	db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	// Get the underlying sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	log.Println("Successfully connected to database")

	if err = db.AutoMigrate(&Message{}, &MessageHistory{}); err != nil {
		return err
	}

	log.Println("Database migration completed successfully")

	return nil
}

func closeDatabase() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
