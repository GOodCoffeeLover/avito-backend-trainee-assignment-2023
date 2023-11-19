.PHONY: docker-compose-up docker-compose-down local-up local-down curls

DOCKER_API_IMAGE="segments/api"
DOCKER_MIGRATION_IMAGE="segments/migration"
KIND_CLUSTER_NAME="local-segments"


docker-compose-up:
	docker compose -f tools/docker-compose.yml up -d --build

docker-compose-down:
	docker compose  -f tools/docker-compose.yml down

local-up:
	docker start postgres || (docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres)
	while ! (pg_isready --host localhost --port 5432); do sleep 1 ;done
	go run ./cmd/api/main.go

local-down:
	docker stop postgres

curls:
	/usr/bin/env bash tools/scripts/curls.sh

k8s-up: setup-kind build-docker-images copy-images-to-kind	
	kubectl get all -A

setup-kind:
	kind create cluster --name ${KIND_CLUSTER_NAME}

copy-images-to-kind: build-docker-images
	kind load docker-image ${DOCKER_MIGRATION_IMAGE} ${DOCKER_API_IMAGE} --name ${KIND_CLUSTER_NAME}

build-docker-images:
	docker build --tag ${DOCKER_API_IMAGE} --file ./tools/Dockerfile .
	docker build --tag ${DOCKER_MIGRATION_IMAGE} --file ./tools/Dockerfile.migration .
