package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var pgdb *gorm.DB


func migrateFromPostgres(sqlitePath string) error {
	if pgdb == nil {
		return nil
	}

	// Check if SQLite is empty
	var sqliteCount int64
	if err := db.Model(&Message{}).Count(&sqliteCount).Error; err != nil {
		return err
	}
	if sqliteCount > 0 {
		log.Println("SQLite already contains data; skipping migration")
		return nil
	}

	// Get counts from Postgres
	var pgMessagesCount, pgHistoryCount int64
	if err := pgdb.Model(&Message{}).Count(&pgMessagesCount).Error; err != nil {
		return err
	}
	if err := pgdb.Model(&MessageHistory{}).Count(&pgHistoryCount).Error; err != nil {
		return err
	}

	if pgMessagesCount == 0 && pgHistoryCount == 0 {
		log.Println("Postgres is also empty; nothing to migrate")
		return nil
	}

	log.Printf("SQLite empty; migrating %d messages and %d history entries from Postgres", pgMessagesCount, pgHistoryCount)

	// Migrate Messages
	var messages []Message
	if pgMessagesCount > 0 {
		if err := pgdb.Find(&messages).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(messages, 100).Error; err != nil {
			return err
		}
	}

	// Migrate MessageHistory
	var history []MessageHistory
	if pgHistoryCount > 0 {
		if err := pgdb.Find(&history).Error; err != nil {
			return err
		}
		if err := db.CreateInBatches(history, 100).Error; err != nil {
			return err
		}
	}

	// Verify row counts match
	var finalMessages, finalHistory int64
	if err := db.Model(&Message{}).Count(&finalMessages).Error; err != nil {
		return err
	}
	if err := db.Model(&MessageHistory{}).Count(&finalHistory).Error; err != nil {
		return err
	}

	if finalMessages != pgMessagesCount || finalHistory != pgHistoryCount {
		log.Printf("Row count mismatch after migration: SQLite (%d/%d) vs Postgres (%d/%d); recreating SQLite and retrying",
			finalMessages, finalHistory, pgMessagesCount, pgHistoryCount)

		// Close SQLite and delete file
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		if err := os.Remove(sqlitePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to remove SQLite file for recreation: %v", err)
			return err
		}

		// Reopen SQLite (this will recreate the file)
		sqliteDSN := sqlitePath
		sep := "?"
		if strings.Contains(sqliteDSN, "?") {
			sep = "&"
		}
		sqliteDSN = sqliteDSN + sep + "_foreign_keys=1&_busy_timeout=5000"
		db, err = gorm.Open(sqlite.Open(sqliteDSN), &gorm.Config{})
		if err != nil {
			return err
		}
		if err := db.AutoMigrate(&Message{}, &MessageHistory{}); err != nil {
			return err
		}

		// Retry migration once
		if pgMessagesCount > 0 {
			if err := db.CreateInBatches(messages, 100).Error; err != nil {
				return err
			}
		}
		if pgHistoryCount > 0 {
			if err := db.CreateInBatches(history, 100).Error; err != nil {
				return err
			}
		}

		// Re-verify
		if err := db.Model(&Message{}).Count(&finalMessages).Error; err != nil {
			return err
		}
		if err := db.Model(&MessageHistory{}).Count(&finalHistory).Error; err != nil {
			return err
		}
		if finalMessages != pgMessagesCount || finalHistory != pgHistoryCount {
			return errors.New("row count still mismatched after SQLite recreation and retry")
		}
	}

	log.Println("Migration from Postgres to SQLite completed successfully")
	return nil
}

func initDatabase(sqlitePath string, postgresURL string) error {
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
	log.Println("SQLite migration completed successfully")

	if postgresURL != "" {
		pgdb, err = gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to Postgres (continuing with SQLite only): %v", err)
			pgdb = nil
			return nil
		}

		pgSQLDB, err := pgdb.DB()
		if err != nil {
			log.Printf("Failed to get Postgres DB handle (continuing with SQLite only): %v", err)
			pgdb = nil
			return nil
		}
		if err = pgSQLDB.Ping(); err != nil {
			log.Printf("Failed to ping Postgres (continuing with SQLite only): %v", err)
			pgdb = nil
			return nil
		}
		log.Println("Successfully connected to Postgres")

		if err = pgdb.AutoMigrate(&Message{}, &MessageHistory{}); err != nil {
			log.Printf("Postgres migration failed (continuing with SQLite only): %v", err)
			pgdb = nil
			return nil
		}
		log.Println("Postgres migration completed successfully")

		// One-time migration from Postgres to SQLite if SQLite is empty
		if err := migrateFromPostgres(sqlitePath); err != nil {
			log.Printf("Migration from Postgres to SQLite failed: %v", err)
			// Continue anyway; SQLite is primary
		}
	}

	return nil
}

func closeDatabase() error {
	var firstErr error

	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			firstErr = err
		} else if err := sqlDB.Close(); err != nil {
			firstErr = err
		}
	}

	if pgdb != nil {
		pgSQLDB, err := pgdb.DB()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
		} else if err := pgSQLDB.Close(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}
