package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	DatabaseURL  string
	SQLitePath   string
}

func loadConfig() *Config {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set; Postgres connection will be skipped")
	}

	sqlitePath := os.Getenv("SQLITE_DB_PATH")
	if sqlitePath == "" {
		if os.Getenv("FLY_APP_NAME") != "" {
			sqlitePath = "/data/eo-bot.sqlite"
		} else {
			sqlitePath = "eo-bot.sqlite"
		}
	}

	return &Config{
		DiscordToken: token,
		DatabaseURL:  dbURL,
		SQLitePath:   sqlitePath,
	}
}
