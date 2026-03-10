build-api:
	cd apps/api && just build

build-web:
	cd apps/web && bun run build

build: build-api build-web

dev-api: build-api
    cd apps/api && ./main

dev-web:
	cd apps/web && bun run dev

dev: dev-web


docker-up:
    systemctl start docker && cd infra && docker-compose up
