SHELL := /bin/bash
GO := docker-compose run --rm go
GO_DETACHED := docker-compose run -d --rm go
GO_FILES := main.go
VERSION := $(shell cat .version)

.PHONY: init build clean run test docker-image

build: src/vendor src/build

clean-build: clean build

init:
	rm -rf src/go.mod src/vendor/
	$(GO) go mod init

src/build:
	$(GO) go build -v -a -installsuffix cgo -o ./build/docker-file
	ls -lach ./src/build/docker-file

src/vendor:
	$(GO) sh -c 'go mod download && go mod vendor -v && ls -lach ./vendor/'

clean:
	rm -rf src/build/

run:
	$(GO) go run .

test:
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file -help
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file dockerfile parse ./test/Dockerfile.valid
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file dockerfile parse ./test/Dockerfile.invalid || true
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file compose parse ./test/docker-compose-valid.yml
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file compose parse ./test/docker-compose-invalid.yml || true
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file env parse ./test/env-valid
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file env parse ./test/env-invalid || true
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file dockerignore parse ./test/ignore-valid
	docker run --rm -v $$(pwd):/srv -w /srv/src alpine:3.9 ./build/docker-file dockerignore check ./test/ignore-valid README.md
