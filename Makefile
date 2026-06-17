.PHONY: dev dev-frontend dev-backend build rebuild down logs ps clean help

# Local
dev:
	docker compose up db backend redis -d

dev-frontend:
	cd frontend && npm run dev

dev-build:
	docker compose up --build db backend redis -d

# Deploy
build:
	cd frontend && npm run build

rebuild:
	docker compose down
	docker builder prune -f
	docker compose up --build backend redis -d

restart:
	docker compose restart backend

# Monitoring
logs:
	docker compose logs backend -f

logs-redis:
	docker compose logs redis -f 

ps:
	docker compose ps

# Database
DB_HOST ?= localhost
db:
	psql -h $(DB_HOST) -U booktracker_user -d booktracker

# Cleanup
down:
	docker compose down

down-volumes:
	docker compose down -v

clean:
	docker builder prune -f 

clean-all:
	docker system prune -f 

# Go
vet:
	cd backend && go vet ./...

tidy:
	cd backend && go mod tidy

# Help
help:
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo '  dev          Start db + backend + redis in background'
	@echo '  dev-build    Same but rebuilds first'
	@echo '  build        Build frontend for production'
	@echo '  rebuild      Full rebuild with cache clear'
	@echo '  restart      Restart backend only'
	@echo '  logs         Follow backend logs'
	@echo '  logs-redis   Follow redis logs'
	@echo '  ps           Show running containers'
	@echo '  down         Stop all containers'
	@echo '  down-volumes Stop and delete volumes'
	@echo '  clean        Clear Docker build cache'
	@echo '  clean-all    Clear all unused Docker resources'
	@echo '  vet          Run go vet on backend'
	@echo '  tidy         Tidy go modules'
	@echo '  db           Connect to RDS via psql'
	@echo ''