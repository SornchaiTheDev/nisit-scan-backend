migrate-up:
	GOOSE_DRIVER=pgx GOOSE_DBSTRING="postgres://nisits-db:bG2LuqBedMjyJ9bVIjxqhQ==@localhost:5432/nisits-scan" GOOSE_MIGRATION_DIR="internal/migrations" goose up

migrate-down:
	GOOSE_DRIVER=pgx GOOSE_DBSTRING="postgres://nisits-db:bG2LuqBedMjyJ9bVIjxqhQ==@localhost:5432/nisits-scan" GOOSE_MIGRATION_DIR="internal/migrations" goose down
