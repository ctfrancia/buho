.PHONY: run-dev

help: ## Show this help
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

confirm: ## Confirm the action
	@read -p "Are you sure? [y/N] " response; \
	if [ "$$response" != "y" ]; then \
		echo "Exiting..."; \
		exit 1; \
	fi

run-dev: ## Run the application in development mode
	@go run ./cmd/...

gen-key: ## show instructions to generate a key
	@echo "To generate a key, run the following command in the root of this project !!!!TESTING ONLY!!!"
	@echo "do not have a password on the key"
	@echo "ssh-keygen -t rsa -b 4096 -f ./id_rsa"
	@echo "copy and paste the public key(id_rsa.pub) to the authorized_keys file in the sftp server"
	@echo "do the same thing in the sftp server and copy the public key to internal/sftp/pub_key file"

setup-db: confirm ## Setup databse locally for testing
	@echo "====REVERTING TO CLEAN STATE===="
	@echo "---Dropping database---"
	psql postgres -c "DROP DATABASE IF EXISTS buho_chess"
	@echo "----------------------"

	@echo "---Dropping user---"
	psql postgres -c "DROP USER IF EXISTS buho_admin"
	@echo "----------------------"

	@echo "---Dropping extensions---"
	psql postgres -c "DROP EXTENSION IF EXISTS citext"
	psql postgres -c 'DROP EXTENSION IF EXISTS "uuid-ossp"'
	@echo "----------------------"

	@echo "====Setting up development environment===="
	@echo "---Creating database---"
	psql postgres -c "CREATE DATABASE buho_chess"
	@echo "----------------------"

	@echo "---creating user and password---"
	psql postgres -c "CREATE USER buho_admin WITH PASSWORD 'pa55word'"
	@echo "----------------------"

	@echo "---changing ownership of the database to user just created---"
	psql postgres -c "ALTER DATABASE buho_chess OWNER TO buho_admin"
	@echo "----------------------"

	@echo "creating extensions for the database with user just created"
	psql chess_admin -d my_chess_website -c "CREATE EXTENSION IF NOT EXISTS citext"
	psql chess_admin -d my_chess_website -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"'
