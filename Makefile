# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify

# ==================================================================================== #
# CSS
# ==================================================================================== #

## css/generate: generate bulma based css
.PHONY: css/generate
css/generate:
	npm run css-build

# ==================================================================================== #
# DB
# ==================================================================================== #

.PHONY: db/createSessionsTable
db/createSessionsTable:
	sqlite3 liberator.db "drop index sessions_expiry_idx" "drop table sessions" "CREATE TABLE sessions (token TEXT PRIMARY KEY,data BLOB NOT NULL,expiry REAL NOT NULL);" "CREATE INDEX sessions_expiry_idx ON sessions(expiry);" ".exit"

.PHONY: db/clearSessions
db/clearSessions:
	sqlite3 liberator.db "DELETE FROM sessions;" ".exit"

# ==================================================================================== #
# App
# ==================================================================================== #

.PHONY: run
run:
	go run ./cmd/web -port=5000