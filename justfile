build-api:
	cd apps/api && just build

build-web:
	cd apps/web && bun run build

build: build-api build-web

dev-api: build-api
    cd apps/api && just dev

dev-web:
	cd apps/web && bun run dev

dev: docker-up

docker-up:
    cd infra && docker-compose up

docker-down:
    cd infra && docker-compose down

