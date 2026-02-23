# ===================================
# LAR CEPPRS - Makefile Profissional
# ===================================

COMPOSE=docker compose
COMPOSE_STAGING=docker compose -f docker-compose.staging.yml
COMPOSE_DEV_INFRA=docker compose -f docker-compose.dev-infra.yml
BACKEND_DIR=backend
FRONTEND_DIR=frontend

.PHONY: help up up-staging down down-staging logs logs-staging backend frontend format clean restart restart-staging status quick-test smoke dev-infra-up dev-infra-down dev-infra-logs dev-infra-status dev-backend dev-frontend dev-seed

help:
	@echo "======================================="
	@echo "PRODUÇÃO (Keycloak externo):"
	@echo " make up           -> Sobe tudo (Docker)"
	@echo " make up-backend   -> Sobe backend + banco"
	@echo " make up-db        -> Sobe apenas o banco de dados"
	@echo " make down         -> Derruba ambiente"
	@echo " make restart      -> Reinicia ambiente"
	@echo ""
	@echo "STAGING (Keycloak local lar-sso):"
	@echo " make up-staging      -> Sobe LAR CEPPRS usando Keycloak do lar-sso"
	@echo " make down-staging    -> Derruba LAR CEPPRS staging"
	@echo " make restart-staging -> Reinicia LAR CEPPRS staging"
	@echo " make logs-staging    -> Logs do staging"
	@echo ""
	@echo "LOCAL (sem Docker):"
	@echo " make backend      -> Roda backend local (binário)"
	@echo " make frontend     -> Roda frontend local (binário)"
	@echo ""
	@echo "DEV HÍBRIDO (recomendado):"
	@echo " make dev-infra-up     -> Sobe infraestrutura local (postgres/mongo/redis/rabbit)"
	@echo " make dev-infra-down   -> Derruba infraestrutura local"
	@echo " make dev-infra-logs   -> Logs da infraestrutura local"
	@echo " make dev-infra-status -> Status da infraestrutura local"
	@echo " make dev-backend      -> Backend local (Go) com infra Docker"
	@echo " make dev-frontend     -> Frontend local (Angular) com infra Docker"
	@echo " make dev-seed         -> Popula professores/alunos para teste local"
	@echo ""
	@echo "OUTROS:"
	@echo " make logs         -> Logs em tempo real (produção)"
	@echo " make logs-backend -> Logs apenas do backend"
	@echo " make logs-db      -> Logs apenas do banco"
	@echo " make quick-test   -> Preflight rápido API/Keycloak"
	@echo " make smoke        -> Smoke RBAC Keycloak (automatizado)"
	@echo " make format       -> Formata código"
	@echo " make clean        -> Limpa docker"
	@echo " make status       -> Status dos containers"
	@echo "======================================="
	@echo ""
	@echo "URLs STAGING:"
	@echo "  LAR CEPPRS Frontend: http://localhost:4201"
	@echo "  LAR CEPPRS Backend:  http://localhost:8081"
	@echo "  Keycloak:       http://localhost:8081 (do lar-sso)"
	@echo "======================================="

up:
	$(COMPOSE) up --build -d

up-backend:
	$(COMPOSE) up --build -d backend

up-db:
	$(COMPOSE) up -d postgres

down:
	$(COMPOSE) down

restart: down up

logs:
	$(COMPOSE) logs -f

logs-backend:
	$(COMPOSE) logs -f backend

logs-db:
	$(COMPOSE) logs -f postgres

backend:
	cd $(BACKEND_DIR) && go run cmd/api/main.go

frontend:
	cd $(FRONTEND_DIR) && npm start

format:
	cd $(BACKEND_DIR) && gofmt -w .
	cd $(FRONTEND_DIR) && npm run format

# --- Staging Commands (usa Keycloak do lar-sso) ---

up-staging:
	@echo "Subindo LAR CEPPRS em modo STAGING (usando Keycloak do lar-sso)..."
	@echo "Certifique-se de que o lar-sso está rodando: cd ../lar-sso && make up"
	$(COMPOSE_STAGING) up --build -d

down-staging:
	$(COMPOSE_STAGING) down

restart-staging: down-staging up-staging

logs-staging:
	$(COMPOSE_STAGING) logs -f

status:
	@echo "=== Containers de Produção ==="
	$(COMPOSE) ps
	@echo ""
	@echo "=== Containers de Staging ==="
	$(COMPOSE_STAGING) ps 2>/dev/null || echo "Staging não está rodando"

clean:
	docker system prune -f

quick-test:
	./scripts/quick_api_test.sh

smoke: quick-test
	./scripts/smoke_rbac_keycloak.sh

# --- Dev híbrido: app local + infra Docker ---
dev-infra-up:
	$(COMPOSE_DEV_INFRA) up -d

dev-infra-down:
	$(COMPOSE_DEV_INFRA) down

dev-infra-logs:
	$(COMPOSE_DEV_INFRA) logs -f

dev-infra-status:
	$(COMPOSE_DEV_INFRA) ps

dev-backend:
	cd $(BACKEND_DIR) && go run cmd/api/main.go

dev-frontend:
	cd $(FRONTEND_DIR) && npm start

dev-seed:
	cd $(BACKEND_DIR) && go run cmd/seed_local/main.go
