
confirm:
	@echo -n 'Are you sure? [y/n] ' && read ans && [ $${ans:-n} = y ]

run/api:
	go run ./cmd/api

db/psql:
	psql ${GOCHAT_DSN}

db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./internal/db/migrations ${name}

db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path=./internal/db/migrations -database=${GOCHAT_DSN} up

db/migrations/down: confirm
	@echo 'Running down migrations by ${num}...'
	migrate -path=./internal/db/migrations -database=${GOCHAT_DSN} down ${num}

.PHONY: run/api db/psql db/migrations/up db/migrations/down
