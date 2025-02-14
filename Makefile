# Configuration
KEY_SIZE ?= 2048
BASE_DIR ?= internal/keys
KEY_TYPES ?= jwt sftp app

# Define standard key names
PRIVATE_KEY = private.pem
PUBLIC_KEY = public.pem

.PHONY: all clean generate-all verify-dirs $(KEY_TYPES) generate-% help confirm run-dev gen-test-keys setup-db

help: ## Show this help
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all-dev-keys: generate-all ## Create all dev keys

generate-all: $(KEY_TYPES) ## Create all key pairs for all types
	@echo "All key pairs generated successfully"

verify-dirs: ## Verify directories exist
	@mkdir -p $(addprefix $(BASE_DIR)/,$(KEY_TYPES))

$(KEY_TYPES): verify-dirs ## Generate keys for each type
	@echo "Generating RSA key pair for $@..."
	@openssl genpkey -algorithm RSA \
		-pkeyopt rsa_keygen_bits:$(KEY_SIZE) \
		-out $(BASE_DIR)/$@/$(PRIVATE_KEY)
	@openssl rsa -pubout \
		-in $(BASE_DIR)/$@/$(PRIVATE_KEY) \
		-out $(BASE_DIR)/$@/$(PUBLIC_KEY)
	@chmod 600 $(BASE_DIR)/$@/$(PRIVATE_KEY)
	@chmod 644 $(BASE_DIR)/$@/$(PUBLIC_KEY)
	@echo "Generated $(BASE_DIR)/$@/$(PRIVATE_KEY)"
	@echo "Generated $(BASE_DIR)/$@/$(PUBLIC_KEY)"

clean: ## Remove all generated keys
	@rm -rf $(BASE_DIR)
	@echo "Removed all generated keys"

generate-%: ## Generate a single key type
	@if echo "$(KEY_TYPES)" | grep -w "$*" > /dev/null; then \
		$(MAKE) $*; \
	else \
		echo "Error: Invalid key type '$*'. Valid types are: $(KEY_TYPES)"; \
		exit 1; \
	fi

confirm: ## Confirm the action
	@read -p "Are you sure? [y/N] " response; \
	if [ "$$response" != "y" ]; then \
		echo "Exiting..."; \
		exit 1; \
	fi

run-dev: ## Run the application in development mode
	@go run ./cmd/...

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
