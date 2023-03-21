TESTS?=$$(go list ./... | egrep -v "mock|docs")
BINARY=go_api_template
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
	docker build -t go_api_template .

docker/build-test:
	docker build -t go_api_template_test . -f Dockerfile.test

docker/run: docker/build
	docker run -e PORT=$(PORT) --rm --name go_api_template  -p $(PORT):$(PORT) go_api_template

.PHONY: docker/clean
docker/clean:
	docker rmi -f go_api_template
	docker rmi -f go_api_template_test
	
docker/test: docker/build-test
	docker run --rm --name go_api_template_test go_api_template_test go test -v $(TESTS) -failfast