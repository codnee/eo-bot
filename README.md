# Discord Bot

A Discord bot built with Go and PostgreSQL

## Project Structure

Simple flat structure with clear separation of concerns:

```
src/
├── main.go      # Application entry point
├── bot.go       # Bot lifecycle management
├── config.go    # Configuration loading
├── database.go  # Database connection & migrations
├── handlers.go  # Discord command handlers
└── models.go    # Database models
```

## Prerequisites

- Go 1.25 or later
- PostgreSQL database
- Discord Bot Token from [Discord Developer Portal](https://discord.com/developers/applications)

## Local Development

1. Copy the example environment file and update with your values:
   ```bash
   cp .env.example .env
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the bot:
   ```bash
   go build -o bin/eo-bot ./src
   ```

4. Run the bot:
   ```bash
   # Using go run
   go run ./src
   
   # Or run the compiled binary
   ./bin/eo-bot
   ```

## Available Commands

- `!ping` - Responds with "Pong! 🏓"
- `!eo` - Fetches and sends a random message from the database

## Database

The bot uses PostgreSQL with auto-migration. Tables are automatically created on startup:
- `messages` - Stores bot messages
- `message_history` - Tracks sent messages (prevents repeating last 20)

To add custom messages to the database:

```sql
INSERT INTO messages (content) VALUES ('Your custom message here');
```
