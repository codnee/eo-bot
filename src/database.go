package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB


func initDatabase(sqlitePath string) error {
	if sqlitePath == "" {
		return errors.New("sqlite path is required")
	}

	if err := os.MkdirAll(filepath.Dir(sqlitePath), 0o755); err != nil && filepath.Dir(sqlitePath) != "." {
		return err
	}

	sqliteDSN := sqlitePath
	sep := "?"
	if strings.Contains(sqliteDSN, "?") {
		sep = "&"
	}
	sqliteDSN = sqliteDSN + sep + "_foreign_keys=1&_busy_timeout=5000"

	var err error
	db, err = gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.Ping(); err != nil {
		return err
	}
	log.Println("Successfully connected to SQLite")

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
