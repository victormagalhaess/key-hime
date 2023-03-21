TESTS?=$$(go list ./... | egrep -v "mock|docs")
BINARY=klggerserver
ENTRY=main.go
PORT=8080

default: build

#host commands

install:
	go get

build: install
	go build -o $(BINARY) $(ENTRY)

run: build
	PORT=$(PORT) ./$(BINARY)

.PHONY: clean
clean:
	go clean
	rm -f $(BINARY) cover.out coverage.html

test:
	go test -v $(TESTS) -failfast -cover 

cover:
	go test -v $(TESTS) -failfast -coverprofile=cover.out
	go tool cover -html=cover.out -o coverage.html

.PHONY: docs
docs:
	echo "Generating docs.."
	swag init

#dockerized commands

docker/build:
	docker build -t klggerserver .

docker/build-test:
	docker build -t klggerserver_test . -f Dockerfile.test

docker/run: docker/build
	docker run -e PORT=$(PORT) --rm --name klggerserver  -p $(PORT):$(PORT) klggerserver

.PHONY: docker/clean
docker/clean:
	docker rmi -f klggerserver
	docker rmi -f klggerserver_test
	
docker/test: docker/build-test
	docker run --rm --name klggerserver_test klggerserver_test go test -v $(TESTS) -failfast