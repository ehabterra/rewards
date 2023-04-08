
.PHONY : generate
generate:
	mkdir -p "internal/pb"
	rm -rf internal/pb/*.go
	protoc --proto_path=proto --go_out=internal/pb --go_opt=paths=source_relative \
        --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
        proto/*.proto

.PHONY: build
build:
	go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/server cmd/server/main.go

.PHONY: up-build
up-build:
	BUILDKIT_PROGRESS=plain docker-compose -f docker/docker-compose.yaml up -d --build

.PHONY: up
up:
	BUILDKIT_PROGRESS=plain docker-compose -f docker/docker-compose.yaml up -d

.PHONY: down
down:
	docker-compose -f docker/docker-compose.yaml down

.PHONY: debug
debug:
	BUILDKIT_PROGRESS=plain docker-compose -f docker/docker-compose.debug.yaml up -d

.PHONY: client
client:
	go run ./cmd/client

.PHONY: migrate
migrate:
	go run ./cmd/migrate

