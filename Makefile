.PHONY: run fmt test cover infra-up infra-up infra-test infra-local infra-local-down infra-upd lint setup testis mocks alltests

# This assumes tflocal is installed https://github.com/localstack/terraform-local

all: infra-down infra-up infra-test
infra-upd:
	cd Docker && docker-compose -f docker-compose.yaml up -d --remove-orphans

infra-up:
	cd Docker && docker-compose -f docker-compose.yaml up --remove-orphans

infra-down:
	cd Docker && docker-compose -f docker-compose.yaml down

infra-test:
	sleep 5 && aws --region eu-east-1 dynamodb list-tables --endpoint-url http://localhost:4566 --no-cli-pager

infra-local:
	cd terraform/infra && \
	export TF_LOG=INFO && \
	tflocal init && \
	tflocal apply -auto-approve

infra-local-down:
	cd terraform/infra && export TF_LOG=INFO && tflocal destroy -auto-approve

fmt:
	go fmt ./... && cd terraform && terraform fmt

lint:
	golangci-lint run

test:
	go test -timeout 50000ms -v ./internal/... -covermode=count -coverprofile=cover.out && go tool cover -func=cover.out


testi:
	go test -timeout 50000ms -v --short ./tests/... -covermode=count -coverprofile=cover.out && go tool cover -func=cover.out

testis:
	go test -timeout 50000ms -v ./tests/... -tags=integration

cover: test
	go tool cover -html=cover.out -o coverage.html

run: fmt lint
	go run ./cmd/simplematchmaking/main.go

setup:  infra-upd infra-local run

docker-build:
	docker buildx build --no-cache --load --platform linux/arm64 -t simplematchmaking --progress plain  .

docker-run:
	docker run --network="host"  --rm  -e MATCHMAKING_CFG=$$MATCHMAKING_CFG -p 4566:4566 -p 8081:8000 simplematchmaking:latest

mocks:
	mockgen -source=internal/core/ports/ports.go -destination=internal/core/ports/mocks/ports_mock.go -package=mocks

alltests: mocks test testi		
