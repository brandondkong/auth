build-api:
	cd apps/api && go build -o api

build-web:
	cd apps/web && bun build

dev-web:
	cd apps/web && bun run dev

dev: dev-web
