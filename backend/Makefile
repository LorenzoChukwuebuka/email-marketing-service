DB_DRIVER := postgres
DB_URL := postgres://postgres:Lorenzo@localhost:5432/email_marketing_service?sslmode=disable
DB_MIGRATIONS_DIR :=./internal/db/migrations
DB_MIGRATE_EXEC := migrate
DB_MIGRATE_ARGS := -database $(DB_DRIVER) -url $(DB_URL) -path $(DB_MIGRATIONS_DIR) -verbose


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

.PHONY: create-migration
create-migration: ## Create a new migration file for a table, e.g., make create-migration table_name=bug_report
ifdef table_name
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -ext sql -dir $(DB_MIGRATIONS_DIR) -seq $(table_name)
else
	@echo "Please provide a table_name argument, e.g., make create-migration table_name=bug_report"
endif

.PHONY: migrate-force
migrate-force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR)  force ${version}


.PHONY: migrate-up
migrate-up: ## Migrate the database schema up to the latest version
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR)  up

.PHONY: migrate-down
migrate-down: ## Rollback the database schema by one migration
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) down


.PHONY: migrate-rollback
migrate-rollback: ## Rollback the database schema by one migration
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) down $(step)


.PHONY: migrate-drop
migrate-drop:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL)  -path=$(DB_MIGRATIONS_DIR) drop -f

.PHONY: migrate-to
migrate-to: ## Migrate the database schema to a specific version, e.g., make migrate-to version=1
ifdef version
	 go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR)  goto $(version)
else
	@echo "Please provide a version argument, e.g., make migrate-to version=1"
endif
