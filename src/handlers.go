package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func fallbackDB(err error) *gorm.DB {
	if err == nil {
		return nil
	}
	if pgdb == nil {
		return nil
	}
	log.Printf("SQLite error; falling back to Postgres: %v", err)
	return pgdb
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if message is a DM and starts with !new
	if m.GuildID == "" && len(m.Content) > 4 && m.Content[:4] == "!new" {
		handleNewMessage(s, m)
		return
	}

	switch m.Content {
	case "!eo":
		handleEo(s, m)
	case "!help":
		handleHelp(s, m)
	}
}

func handleNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Extract and trim the message content after "!new"
	content := strings.TrimSpace(m.Content[4:])
	if content == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Please provide a message after !new")
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Create and save the new message
	message := Message{
		Content: content,
	}
	result := db.Create(&message)
	if result.Error != nil {
		if fb := fallbackDB(result.Error); fb != nil {
			result = fb.Create(&message)
		}
	}

	if result.Error != nil {
		log.Printf("Error saving message to database: %v", result.Error)
		_, err := s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't save your message. Please try again later.")
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Send confirmation
	_, err := s.ChannelMessageSend(m.ChannelID, "✅ Message saved successfully!")
	if err != nil {
		log.Printf("Error sending confirmation: %v", err)
	}
}

func handleEo(s *discordgo.Session, m *discordgo.MessageCreate) {
	activeDB := db

	// Get the last 50 message IDs from history for THIS channel
	var recentMessageIDs []uint
	result := activeDB.Model(&MessageHistory{}).
		Where("channel_id = ?", m.ChannelID).
		Order("sent_at DESC").
		Limit(50).
		Pluck("message_id", &recentMessageIDs)
	if result.Error != nil {
		if fb := fallbackDB(result.Error); fb != nil {
			activeDB = fb
			recentMessageIDs = nil
			activeDB.Model(&MessageHistory{}).
				Where("channel_id = ?", m.ChannelID).
				Order("sent_at DESC").
				Limit(50).
				Pluck("message_id", &recentMessageIDs)
		}
	}

	// Select a random message NOT in the recent list
	var message Message
	query := activeDB.Order("RANDOM()")
	if len(recentMessageIDs) > 0 {
		query = query.Where("id NOT IN ?", recentMessageIDs)
	}
	result = query.First(&message)
	if result.Error != nil && activeDB == db {
		if fb := fallbackDB(result.Error); fb != nil {
			activeDB = fb
			query = activeDB.Order("RANDOM()")
			if len(recentMessageIDs) > 0 {
				query = query.Where("id NOT IN ?", recentMessageIDs)
			}
			result = query.First(&message)
		}
	}

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
	if err := activeDB.Create(&history).Error; err != nil {
		log.Printf("Error recording message history: %v", err)
		if activeDB == db {
			if fb := fallbackDB(err); fb != nil {
				_ = fb.Create(&history).Error
			}
		}
	}

	// Send the message to Discord
	_, err := s.ChannelMessageSend(m.ChannelID, `"` + message.Content + `"`)

	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := `📝 **Available Commands** 📝

- ` + "`!new <message>`" + ` - Save a new message to the database (DMs only)
- ` + "`!eo`" + ` - Get a random message from the database
- ` + "`!help`" + ` - Show this help message

For support, contact the bot administrator.`

	_, err := s.ChannelMessageSend(m.ChannelID, helpMessage)
	if err != nil {
		log.Printf("Error sending help message: %v", err)
	}
}