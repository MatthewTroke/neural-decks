rebuild-dev: down-dev build-dev up-dev
all: down-dev build-dev up-dev

build-dev:
	docker compose -f docker-compose.dev.yml build $(N) --no-cache

up-dev:
	docker compose -f docker-compose.dev.yml up $(N) -d

down-dev:
	docker compose -f docker-compose.dev.yml down $(N) || true

logs-dev:
	docker compose -f docker-compose.dev.yml logs -f backend