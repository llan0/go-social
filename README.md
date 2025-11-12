# go-social

Using,
- Go 1.25
- Docker & Docker Compose (for Postgres)
- Chi router: https://github.com/go-chi/chi
- Migrate CLI for running migrations: https://github.com/golang-migrate/migrat
- Air for live reloading: https://github.com/air-verse/air
- Go swagger for Docs: https://github.com/swaggo/swag
- Frontend (TDB - but Svelte or NextJS)

Quick setup

1. Start the DB
   ```sh
   docker-compose up -d
   ```

2. Configure the DB connection
   - Create a `.envrc` or `.env` with your DB address. Example (used by Makefile/migrate tools):
     ```
     DB_ADDR=postgres://admin:adminpassword@localhost:5432/social?sslmode=disable
     ```
   - If using direnv, run `direnv allow` after creating `.envrc`.

3. Run migrations
   - Using the `migrate` CLI:
     ```sh
     migrate -path=./cmd/migrate/migrations -database="${DB_ADDR}" up
     ```
   - Or use the Makefile:
     ```sh
     make migrate-up
     ```

4. Build and run the API
   - Build:
     ```sh
     go build -o ./bin/main ./cmd/api/*.go
     ```
   - Run:
     ```sh
     ./bin/main
     ```
   - OR just `air`

Useful commands
- Seed DB: `make seed`
- Recreate migrations: `make migration NAME`
- Run migrations down: `make migrate-down NAME` (see Makefile)
