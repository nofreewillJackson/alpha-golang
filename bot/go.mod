module github.com/nofreewilljackson/alpha-golang/bot

go 1.23.1

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/jackc/pgx/v4 v4.18.3
	github.com/joho/godotenv v1.5.1
	github.com/nofreewilljackson/alpha-golang/common v0.0.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sashabaranov/go-openai v1.29.2
// other dependencies...
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/nofreewilljackson/alpha-golang/common => ../common
