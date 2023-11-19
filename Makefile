.PHONY: docker-compose-up docker-compose-down local-up local-down curls

DOCKER_API_IMAGE="segments/api"
DOCKER_MIGRATION_IMAGE="segments/migration"
KIND_CLUSTER_NAME="segments-local"


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

k8s-curls:
	SEGMETNS_API_ADDR="--resolve segments-api.local:80:$$(docker network inspect kind | jq -r ".[].IPAM.Config[0].Gateway") segments-api.local" /usr/bin/env bash tools/scripts/curls.sh

k8s-down:
	kind delete cluster --name ${KIND_CLUSTER_NAME}

k8s-up: setup-kind copy-images-to-kind
	helm dependency build .helm
	helm upgrade -i segments-api .helm --namespace segments-api --create-namespace --atomic --debug
	echo "to access app run 'curl -sL --resolve segments-api.local:80:$$(docker network inspect kind | jq -r ".[].IPAM.Config[0].Gateway") segments-api.local/v1/user'"

setup-kind:
	(kind get clusters | grep ${KIND_CLUSTER_NAME}) || kind create cluster --name ${KIND_CLUSTER_NAME} --config=./tools/kind-cluster.yml
	/usr/bin/env bash ./tools/scripts/install-ingress.sh

copy-images-to-kind: build-docker-images
	kind load docker-image ${DOCKER_MIGRATION_IMAGE} ${DOCKER_API_IMAGE} --name ${KIND_CLUSTER_NAME}

build-docker-images:
	docker build --tag ${DOCKER_API_IMAGE} --file ./tools/Dockerfile .
	docker build --tag ${DOCKER_MIGRATION_IMAGE} --file ./tools/Dockerfile.migration .

