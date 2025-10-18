# EO Bot

A Discord bot that responds with random messages from a predefined set of responses.

## Project Structure

```
src/
├── main.go      # Application entry point
├── bot.go       # Bot initialization and connection
├── config.go    # Environment configuration
├── database.go  # Database connection and queries
├── handlers.go  # Command handlers
└── models.go    # Data models
```

## Prerequisites

- Go 1.25 or later
- PostgreSQL database
- Discord Bot Token

## Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
   Update the values in `.env` with your configuration.

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the bot:
   ```bash
   go run ./src
   ```

## Usage

The bot responds to the following commands:

- `!eo` - Get a random message that hasn't been sent recently in the current channel
- `!new [message]` - Add a new message to the bot's database (DM only)
- `!help` - Show available commands

### Adding New Messages

To add new messages, send the bot a direct message (DM) with the following format:
```
!new Your message here
```

### Message Rotation
- The bot keeps track of the last 20 messages sent in each channel
- It will avoid repeating these messages in the same channel
- Messages are selected randomly from the remaining pool

## Deployment

See `QUICK_START.md` for deployment instructions.
