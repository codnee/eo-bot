package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
	case "!ping":
		handlePing(s, m)
	case "!eo":
		handleEo(s, m)
	case "!crawl":
		handleCrawl(s, m)
	}
}

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong! 🏓 - New EO Bot coming soon")
	if err != nil {
		log.Printf("Error sending message: %v", err)
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

func handleCrawl(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Send initial message
    msg, err := s.ChannelMessageSend(m.ChannelID, "🔄 Crawling channel history for EO Bot messages...")
    if err != nil {
        log.Printf("Error sending initial message: %v", err)
        return
    }

    var messagesSaved, messagesSkipped int
    var beforeID string

    // Keep fetching messages until we reach the beginning of the channel
    for {
        // Fetch messages (100 at a time, which is the maximum allowed by Discord)
        messages, err := s.ChannelMessages(m.ChannelID, 100, beforeID, "", "")
        if err != nil {
            log.Printf("Error fetching messages: %v", err)
            s.ChannelMessageEdit(m.ChannelID, msg.ID, "❌ Error: Could not fetch messages from this channel")
            return
        }

        // If no more messages, we're done
        if len(messages) == 0 {
            break
        }

        // Process messages
        for _, message := range messages {
            // Skip if not from a bot or not from "EO Bot"
            if !message.Author.Bot || message.Author.Username != "EO Bot" {
                continue
            }

            // Skip if message is empty or is a command
            content := strings.TrimSpace(message.Content)
            if content == "" || strings.HasPrefix(content, "!") {
                continue
            }

            // Try to save the message (will skip duplicates due to unique constraint)
            msg := Message{
                Content: content,
            }
            result := db.Create(&msg)

            if result.Error != nil {
                // Skip duplicate messages (error 19 is SQLite's constraint violation)
                if !strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
                    log.Printf("Error saving message: %v", result.Error)
                }
                messagesSkipped++
            } else {
                messagesSaved++

                // Add to message history
                history := MessageHistory{
                    MessageID: msg.ID,
                    ChannelID: m.ChannelID,
                }
                if err := db.Create(&history).Error; err != nil {
                    log.Printf("Error saving message history: %v", err)
                }
            }
        }

        // Update the beforeID to the ID of the last message in this batch
        beforeID = messages[len(messages)-1].ID

        // Update status message
        s.ChannelMessageEdit(m.ChannelID, msg.ID, 
            fmt.Sprintf("🔄 Crawling... Found %d messages (%d saved, %d skipped)", 
                messagesSaved + messagesSkipped, messagesSaved, messagesSkipped))

        // Small delay to avoid rate limiting
        time.Sleep(500 * time.Millisecond)
    }

    // Send final message
    s.ChannelMessageEdit(m.ChannelID, msg.ID, 
        fmt.Sprintf("✅ Crawl complete! Processed %d messages.\n- Saved: %d\n- Skipped (duplicates): %d",
            messagesSaved + messagesSkipped, messagesSaved, messagesSkipped))
}