BINARY_NAME=Tasker

all: build run tests build-docker run-docker

build:
	go build -o ${BINARY_NAME} cmd/main.go

run:
	go run cmd/api/main.go

tests:
	go test -v ./test/...

build-docker:
	docker build -t ${BINARY_NAME}:v1 .

run-docker:
	docker run --rm -d -p 9090:9090/tcp ${BINARY_NAME}:v1