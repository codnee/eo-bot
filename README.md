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

### Running the Bot

#### Option 1: Quick Run (Development)
```bash
go run ./src
```

#### Option 2: Build and Run (Production-like)
```bash
# Build
go build -o bin/eo-bot ./src

# Run
./bin/eo-bot
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


### Adding a New Command

1. Open `src/handlers.go`
2. Add your handler function:
   ```go
   func handleYourCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
       s.ChannelMessageSend(m.ChannelID, "Your response")
   }
   ```
3. Register it in the `messageCreate` switch statement:
   ```go
   case "!yourcommand":
       handleYourCommand(s, m)
   ```

### Common Issues

**"DISCORD_BOT_TOKEN not found"**
- Make sure `.env` file exists
- Check that `DISCORD_BOT_TOKEN=your_token` is set

**"DATABASE_URL not found"**
- Make sure `.env` file has `DATABASE_URL`
- Format: `postgres://user:password@host:port/dbname?sslmode=disable`

**"No messages in database"**
- Manually insert: `psql $DATABASE_URL -c "INSERT INTO messages (content) VALUES ('Test');"`
