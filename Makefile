.DEFAULT_GOAL = all

CMDDIR=cmd/toysapiserver/main.go

# Detect the os for naming reasons
OS := $(shell uname -s | awk '{print tolower($$0)}')

# Get a short hash of the git had for building images
TAG = $$(git rev-parse --short HEAD)

# Name of actual binary to create
BINARY = otto-rest-api

# GOARCH tells go build which arch. to use
GOARCH = amd64

.PHONY: all
all: test build docker

 # Run the application after building it first
.PHONY: run
run: build
	./$(BINARY)-$(OS)-$(GOARCH)

# Build simply builds the application
.PHONY: build
build:
	env CGO_ENABLED=0 GOOS=$(OS) GOARCH=${GOARCH} go build -o ${BINARY}-$(OS)-${GOARCH} $(CMDDIR) ;

# Build the docker image
.PHONY: docker
docker:
	docker build -t lucku/$(BINARY):$(GOARCH)-$(TAG) .

# Run the Docker image (will only properly work on Linux)
.PHONY: docker-run
docker-run: docker
	docker run --rm --network host --name $(BINARY) lucku/$(BINARY):$(GOARCH)-$(TAG)

# Run unit tests
.PHONY: test
test:
	go test -v ./...

# Generate a coverage report
.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Remove coverage report, binary, and Docker image
.SILENT: clean
.PHONY: clean
clean:
	go clean $(CMDDIR)
	@rm -f ${BINARY}-$(OS)-${GOARCH}
	@rm -f coverage.out
	docker rmi -f lucku/$(BINARY):$(GOARCH)-$(TAG) > /dev/null 2>&1