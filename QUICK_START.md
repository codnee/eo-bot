# Quick Start Guide

## Development Workflow

### Initial Setup
```bash
# 1. Clone and enter the directory
cd eo-bot

# 2. Copy environment file
cp .env.example .env

# 3. Edit .env with your credentials
# Add your DISCORD_BOT_TOKEN and DATABASE_URL

# 4. Install dependencies
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

### File Structure Quick Reference

| File | Purpose | Key Functions/Types |
|------|---------|---------------------|
| `main.go` | Entry point | `main()` |
| `config.go` | Load env variables | `loadConfig()` |
| `database.go` | Database connection | `initDatabase()`, global `db` |
| `models.go` | Data structures | `Message` |
| `handlers.go` | Command logic | `handlePing()`, `handleEo()` |
| `bot.go` | Bot lifecycle | `newBot()`, `start()`, `stop()` |

### Testing Commands

Once the bot is running, test in Discord:
- `!ping` - Should respond with "Pong! 🏓"
- `!eo` - Should return a random message from the database

### Common Issues

**"DISCORD_BOT_TOKEN not found"**
- Make sure `.env` file exists
- Check that `DISCORD_BOT_TOKEN=your_token` is set

**"DATABASE_URL not found"**
- Make sure `.env` file has `DATABASE_URL`
- Format: `postgres://user:password@host:port/dbname?sslmode=disable`

**"No messages in database"**
- Manually insert: `psql $DATABASE_URL -c "INSERT INTO messages (content) VALUES ('Test');"`

### Building for Production

```bash
# Build binary
CGO_ENABLED=1 go build -o bin/eo-bot ./src

# Or use Docker
docker build -t eo-bot .
docker run --env-file .env eo-bot
```

## Next Steps

- Read [README.md](README.md) for deployment instructions
