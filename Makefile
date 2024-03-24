setup:
	@echo "=== Installing Tools ==="
	@bash ./scripts/install_tools.sh
	@echo "=== Tools Installed ==="

lint:
	@echo "=== Running Linter ==="
	@./bin/golangci-lint run
	@echo "=== Linter Completed ==="

air:
	@echo "=== Running AIR ==="
	@./bin/air
	@echo "=== AUR Completed ==="

compose:
	@echo "=== Starting Server ==="
	@docker compose -f docker-compose.yml up --build
	@echo "=== Development Server Stopped ==="
