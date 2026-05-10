include .envrc


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo "Usage: make [target] [args]"
	@sed -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n]' && read ans && [ $${ans:-n} = y ]

# ================================================================================== #
# DEVELOPMENT
# ================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api/ -db-dsn $$DSN

## db/migraation/new name=$1: create a new up/down migration files with the given name
.PHONY: db/migration/new
db/migration/new:
	migrate create -seq -ext sql -dir ./migrations/ -seq ${name}

## db/migration/up: run all up migrations
.PHONY: db/migration/up
db/migration/up: confirm
	@migrate -database $$DSN -path ./migrations/ up ${steps}

## db/migration/down: run the last down migration
.PHONY: db/migration/down
db/migration/down:
	@migrate -database $$DSN -path ./migrations/ down ${steps}
# ================================================================================== #
# BUILD
# ================================================================================== #

## build/api: build the cmd/api application
.PHONY:
build/api:
	@echo "Building cmd/api application..."
	@go build -o bin/api ./cmd/api
	




