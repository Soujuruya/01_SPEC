# Makefile
ENV_FILE := ./config/.env

DC := docker compose --env-file $(ENV_FILE)

DOCKER_WAIT := 60

.PHONY: up stub stop clean

up: stub docker

stub:
	@echo "[Makefile] Starting webhook stub on port 9090..."
	@nohup go run ./webhook_stub/main.go > webhook_stub.log 2>&1 & echo $$! > .webhook_stub.pid
	@sleep 1
	@echo "[Makefile] Webhook stub started with PID $$(cat .webhook_stub.pid)"do

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
