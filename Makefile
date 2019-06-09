.DEFAULT_GOAL = all

# Go related commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test ./...
CMDDIR=cmd/toysapiserver/main.go

# Detect the os so that we can build proper statically linked binary
OS := $(shell uname -s | awk '{print tolower($$0)}')

# Get a short hash of the git had for building images
TAG = $$(git rev-parse --short HEAD)

# Name of actual binary to create
BINARY = otto-rest-api

# GOARCH tells go build which arch. to use
GOARCH = amd64

.PHONY: all
all: test build docker

 # Runs the application after building it first
.PHONY: run
run: build
	./$(BINARY)-$(OS)-$(GOARCH)

# Build simply builds the application
.PHONY: build
build:
	env CGO_ENABLED=0 GOOS=$(OS) GOARCH=${GOARCH} go build -o ${BINARY}-$(OS)-${GOARCH} $(CMDDIR) ;

# Docker build internally (within Dockerfile) triggers "make bin", which creates a "linux" binary.
.PHONY: docker
docker:
	docker build -t lucku/$(BINARY):$(GOARCH)-$(TAG) .

# Runs unit tests.
.PHONY: test
test:
	$(GOTEST)

# Generates a coverage report
.PHONY: cover
cover:
	${GOCMD} test -coverprofile=coverage.out ./... && ${GOCMD} tool cover -html=coverage.out

# Remove coverage report and the binary.
.SILENT: clean
.PHONY: clean
clean:
	$(GOCLEAN) $(CMDDIR)
	@rm -f ${BINARY}-$(OS)-${GOARCH}
	@rm -f coverage.out