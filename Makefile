include .env
export

MIGRATE_IMAGE := ghcr.io/golang-migrate/migrate:latest
DB_URL        := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATE_CMD   := docker run --rm --network host \
                   -v $(PWD)/migrations:/migrations \
                   $(MIGRATE_IMAGE) \
                   -path /migrations -database "$(DB_URL)"

# ─── Migration ───────────────────────────────────────────────────────────────

.PHONY: migrate-up
migrate-up:
	$(MIGRATE_CMD) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_CMD) down 1

.PHONY: migrate-down-all
migrate-down-all:
	$(MIGRATE_CMD) down -all

.PHONY: migrate-version
migrate-version:
	$(MIGRATE_CMD) version

.PHONY: migrate-create
migrate-create:
	@[ -n "$(name)" ] || (echo "Usage: make migrate-create name=<migration_name>"; exit 1)
	docker run --rm -v $(PWD)/migrations:/migrations $(MIGRATE_IMAGE) \
		create -ext sql -dir /migrations -seq $(name)

# ─── Docker Compose ──────────────────────────────────────────────────────────

.PHONY: up
up:
	docker compose up --build

.PHONY: up-d
up-d:
	docker compose up -d --build

.PHONY: down
down:
	docker compose down

.PHONY: down-v
down-v:
	docker compose down -v

# ─── Dev ─────────────────────────────────────────────────────────────────────

.PHONY: run
run:
	go run ./cmd/server

.PHONY: build
build:
	go build -o bin/server ./cmd/server

.PHONY: test
test:
	go test ./...

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: lint
lint:
	golangci-lint run
