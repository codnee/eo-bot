package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := loadConfig()

	if err := initDatabase(cfg.SQLitePath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer closeDatabase()

	discordBot, err := newBot(cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord bot: %v", err)
	}

	if err := discordBot.start(); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
	defer discordBot.stop()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down gracefully...")
}
