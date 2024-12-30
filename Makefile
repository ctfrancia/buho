.PHONY: run-dev

help: ## Show this help
	@echo "Usage: make [target]"
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run-dev: ## Run the application in development mode
	@go run ./cmd/...

gen-key: ## show instructions to generate a key
	@echo "To generate a key, run the following command in the root of this project !!!!TESTING ONLY!!!"
	@echo "do not have a password on the key"
	@echo "ssh-keygen -t rsa -b 4096 -f ./id_rsa"
	@echo "copy and paste the public key(id_rsa.pub) to the authorized_keys file in the sftp server"
	@echo "do the same thing in the sftp server and copy the public key to internal/sftp/pub_key file"
