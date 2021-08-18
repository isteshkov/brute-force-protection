.PHONY: lint
lint:
	go get -u golang.org/x/lint/golint
	golangci-lint run

.PHONY: run
run:
	sudo docker-compose up --build

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: build
build:
	go build -o brute_force_protection .

.PHONY: test
test:
	go test -race -count 100 ./...

.PHONY: proto
proto:
	protoc --go_out=./contract/ --go-grpc_out=./contract/ --go-grpc_opt=require_unimplemented_servers=false --experimental_allow_proto3_optional ./contract/*.proto
