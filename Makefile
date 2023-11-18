.PHONY: docker-compose-run docker-compose-down local-run local-down

docker-compose-run:
	docker compose -f tools/docker-compose.yml up -d --build

docker-compose-down:
	docker compose  -f tools/docker-compose.yml down
local-run:
	docker start postgres || (docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres && sleep 5)
	go run ./cmd/api/main.go

local-down:
	docker stop postgres