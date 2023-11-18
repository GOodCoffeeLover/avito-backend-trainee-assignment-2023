.PHONY: docker-compose-up docker-compose-down local-up local-down curls

docker-compose-up:
	docker compose -f tools/docker-compose.yml up -d --build

docker-compose-down:
	docker compose  -f tools/docker-compose.yml down
local-up:
	docker start postgres || (docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres && sleep 5)
	go run ./cmd/api/main.go

local-down:
	docker stop postgres

curls:
	/usr/bin/env bash tools/scripts/curls.sh