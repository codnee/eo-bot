package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	SQLitePath   string
}

func loadConfig() *Config {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
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
		SQLitePath:   sqlitePath,
	}
}
