package main

import (
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Message) TableName() string {
	return "messages"
}

type MessageHistory struct {
	ID        uint      `gorm:"primaryKey"`
	MessageID uint      `gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Message   Message   `gorm:"foreignKey:MessageID"`
	ChannelID string    `gorm:"not null;index"`
	SentAt    time.Time `gorm:"autoCreateTime;index"`
}

func (MessageHistory) TableName() string {
	return "message_history"
}
