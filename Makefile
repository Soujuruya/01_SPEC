ENV_FILE := ./config/.env
LOCAL_CONFIG := ./config/local.env
MIGRATIONS_PATH := ./migrations

DC := docker compose --env-file $(ENV_FILE)

.PHONY: up stub docker stop clean \
        migrate-up migrate-down migrate-version

up: stub docker

stub:
	@echo "[Makefile] Starting webhook stub on port 9090..."
	@nohup go run ./webhook_stub/main.go > webhook_stub.log 2>&1 & echo $$! > .webhook_stub.pid
	@sleep 1
	@echo "[Makefile] Webhook stub started with PID $$(cat .webhook_stub.pid)"

docker:
	@echo "[Makefile] Starting Docker services..."
	$(DC) up --build -d

stop:
	@echo "[Makefile] Stopping webhook stub..."
	@if [ -f .webhook_stub.pid ]; then \
		PID=$$(cat .webhook_stub.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			kill $$PID && echo "[Makefile] Stub stopped"; \
		else \
			echo "[Makefile] No stub process found"; \
		fi; \
		rm -f .webhook_stub.pid; \
	fi
	@echo "[Makefile] Stopping Docker services..."
	$(DC) down

clean: stop
	@echo "[Makefile] Removing Docker volumes..."
	$(DC) down -v --remove-orphans

migrate-up:
	go run ./cmd/migrate/main.go \
		-config=$(LOCAL_CONFIG) \
		-path=$(MIGRATIONS_PATH) \
		-command=up

migrate-down:
	go run ./cmd/migrate/main.go \
		-config=$(LOCAL_CONFIG) \
		-path=$(MIGRATIONS_PATH) \
		-command=down

migrate-version:
	go run ./cmd/migrate/main.go \
		-config=$(LOCAL_CONFIG) \
		-path=$(MIGRATIONS_PATH) \
		-command=version
