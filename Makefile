.PHONY: run
run:


.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: build
build:
	go build -o brute_force_protection .

.PHONY: tests
tests:


.PHONY: proto
proto:
	protoc --go_out=./contract/ --go-grpc_out=./contract/ --go-grpc_opt=require_unimplemented_servers=false --experimental_allow_proto3_optional ./contract/*.proto
