package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "!ping":
		handlePing(s, m)
	case "!eo":
		handleEo(s, m)
	}
}

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong! 🏓 - New EO Bot coming soon")
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleEo(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get the last 20 message IDs from history for THIS channel
	var recentMessageIDs []uint
	db.Model(&MessageHistory{}).
		Where("channel_id = ?", m.ChannelID).
		Order("sent_at DESC").
		Limit(20).
		Pluck("message_id", &recentMessageIDs)

	// Select a random message NOT in the recent list
	var message Message
	query := db.Order("RANDOM()")
	if len(recentMessageIDs) > 0 {
		query = query.Where("id NOT IN ?", recentMessageIDs)
	}
	result := query.First(&message)

	if result.Error != nil {
		log.Printf("Error fetching message from database: %v", result.Error)
		_, err := s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't fetch a message from the database. 😔")
		if err != nil {
			log.Printf("Error sending error message: %v", err)
		}
		return
	}

	// Record this message in history with channel ID
	history := MessageHistory{
		MessageID: message.ID,
		ChannelID: m.ChannelID,
	}
	if err := db.Create(&history).Error; err != nil {
		log.Printf("Error recording message history: %v", err)
	}

	// Send the message to Discord
	_, err := s.ChannelMessageSend(m.ChannelID, message.Content)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
