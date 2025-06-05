# Gopher Status

My first project using Golang. It's a small chaotic server that can check availability of other services 

## Installation

1. Clone repo
2. Create `.env` file
3. Provide data:

```
TELEGRAM_BOT_TOKEN=<>

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PSW=postgres
DB_NAME=gopher_status

JWT_SECRET=<>
```

4. Run `go run ./cmd/main.go`

## Usage

You can add new services and watch them using web form, or ask telegram bot via `/status` command to give you info on your services
