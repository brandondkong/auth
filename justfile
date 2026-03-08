build-api:
	cd apps/api && just build

build-web:
	cd apps/web && bun build

dev-api: build-api
    cd apps/api && ./main

dev-web:
	cd apps/web && bun run dev

dev: dev-web
