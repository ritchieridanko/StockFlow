COMPOSE_FILE := compose.yaml

help:
	@echo "Available commands:"
	@echo " make docker-build                 Build the services"
	@echo " make docker-up                    Run up the services"
	@echo " make docker-down                  Drop the services"
	@echo " make docker-migrator-up           Apply all database up-migrations"
	@echo " make docker-migrator-down         Apply all database down-migrations"
	@echo " make docker-start                 Build and run up the services, and apply up-migrations"

# ---------- DOCKER COMMANDS ----------
docker-build:
	docker compose -f $(COMPOSE_FILE) build

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down -v

docker-migrator-up:
	docker compose run --rm stockflow-migrator -up

docker-migrator-down:
	docker compose run --rm stockflow-migrator -down 0

docker-start:
	make docker-build
	make docker-up
	make docker-migrator-up
