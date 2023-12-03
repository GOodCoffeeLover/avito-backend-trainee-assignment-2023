.PHONY: docker-compose-up docker-compose-down local-up local-down curls

DOCKER_API_IMAGE="segments/api"
DOCKER_MIGRATION_IMAGE="segments/migration"
DOCKER_IMAGE_TAG="$$(git rev-parse --short HEAD)"
KIND_CLUSTER_NAME="segments-local"

tag: 
	echo ${DOCKER_IMAGE_TAG}

docker-compose-up:
	docker compose -f tools/docker-compose.yml up -d --build

docker-compose-down:
	docker compose  -f tools/docker-compose.yml down

local-up:
	docker start postgres || (docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 --rm postgres)
	while ! (pg_isready --host localhost --port 5432); do sleep 1 ;done
	go run ./tools/migration/migration.go
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
	helm upgrade -i segments-api .helm --namespace segments-api --create-namespace --atomic --debug --set "image.tag=${DOCKER_IMAGE_TAG}" --set "migration.image.tag=${DOCKER_IMAGE_TAG}"
	kubectl config set-context --current --namespace=segments-api
	echo "\nTo access app on linux run 'curl -sL --resolve segments-api.local:80:$$(docker network inspect kind | jq -r ".[].IPAM.Config[0].Gateway") segments-api.local/v1/user' or 'make k8s-curls'\n"\
	"Or to run comman inside cluster 'kubectl run -i --rm --quiet=true wget-users --image=busybox --restart=Never -- sh -c "wget -qO- segments-api.local/v1/user | xargs echo" 2>/dev/null'"

setup-kind:
	(kind get clusters | grep ${KIND_CLUSTER_NAME}) || kind create cluster --name ${KIND_CLUSTER_NAME} --config=./tools/kind-cluster.yml
	kubectl config set-context kind-"${KIND_CLUSTER_NAME}" 
	/usr/bin/env bash ./tools/scripts/install-ingress.sh

copy-images-to-kind: build-docker-images
	kind load docker-image "${DOCKER_API_IMAGE}":"${DOCKER_IMAGE_TAG}" "${DOCKER_MIGRATION_IMAGE}":"${DOCKER_IMAGE_TAG}" --name ${KIND_CLUSTER_NAME}

build-docker-images:
	docker build --tag "${DOCKER_API_IMAGE}":"${DOCKER_IMAGE_TAG}" --file ./tools/Dockerfile .
	docker build --tag "${DOCKER_MIGRATION_IMAGE}":"${DOCKER_IMAGE_TAG}" --file ./tools/Dockerfile.migration .

