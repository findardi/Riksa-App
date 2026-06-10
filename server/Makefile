include configs/.env
export

MIGRATIONS_DIR := migrations
DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: migrate-create migrate-up migrate-down migrate-status sqlc

## migrate-create name=create_users
migrate-create:
	goose -dir $(MIGRATIONS_DIR) -s create $(name) sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

sqlc:
	sqlc generate -f configs/sqlc.yaml
