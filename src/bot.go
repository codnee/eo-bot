package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func newBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.AddHandler(messageCreate)

	return &Bot{
		Session: session,
	}, nil
}

func (b *Bot) start() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	return nil
}

func (b *Bot) stop() error {
	return b.Session.Close()
}
