build-api:
	go build  -o  service-api ./api

run-api:
	go run ./api/main.go

build-web:
	go build  -o  service-web ./web

run-web:
	go run ./web/main.go
build-worker:
	go build  -o  service-worker ./worker

run-worker:
	go run ./worker/main.go


.PHONY: run-api build-api run-web build-web run-worker build-worker

.DEFAULT_GOAL := run-api