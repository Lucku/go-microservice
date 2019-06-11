# Otto Coding Challenge

This project is the implementation of a simple RESTful microservice, which queries an existing API and has the functionality to filter and transform its contents.

## Running the application

Run `make all` to execute all the tests, builds the executable and build the docker image. If only the docker image is needed, only run `make docker` or `docker build -t docker <image-name> .`. To just build and run the application, call `make run`, or `make docker-run` respectively for building and running the docker container. The latter command will only properly work on Linux hosts, since it lets the Docker container use the host network driver to run directly on the port which is specified inside the server. On other hosts, the APIs port has to be properly mapped.

## Configuration

There is a configuration file `config.yaml` in the project root which allows to configure the server. This file has to reside next to the server binary. It allows to configure the follow parameters:

- apiKey: The API key of the queried API (default: key given in task description)
- port: The port that the server will listen on (default: 8080)
- logging-level: Level of logging in the application (debug|info|warn|error|fatal) (default: info)