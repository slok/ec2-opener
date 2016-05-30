.PHONY: default deps  base build dev shell start stop image push
SHELL := $(shell which bash)
DOCKER := $(shell command -v docker)
DOCKER_COMPOSE := $(shell command -v docker-compose)
SERVICE_NAME := ec2-open-port
IMAGE_NAME := slok/$(SERVICE_NAME)
GID := $(shell id -g)
UID := $(shell id -u)
VERSION ?= $(shell cat VERSION)

default: build

deps:
ifndef DOCKER
  $(error "Docker is not available. Please install docker")
endif
ifndef DOCKER_COMPOSE
  $(error "docker-compose is not available. Please install docker-compose")
endif

base: deps
	docker build --build-arg gid=$(GID) --build-arg uid=$(UID) -t $(IMAGE_NAME)_base:latest .

build: base
	cd environment/dev && docker-compose build

dev: build
	cd environment/dev && \
	( docker-compose run --rm $(SERVICE_NAME) bash -c "go run *.go"; \
		docker-compose stop; \
		docker-compose rm -f; )

shell: build
	cd environment/dev && docker-compose run --rm $(SERVICE_NAME) /bin/bash


start: build
	cd environment/dev && \
		docker-compose up -d

stop:
	cd environment/dev && ( \
		docker-compose stop; \
		docker-compose rm -f; \
		)

#image: base
#	docker build -t $(SERVICE_NAME) -t $(REPOSITORY):$(VERSION) -t $(REPOSITORY):latest -f environment/prod/Dockerfile .
#
#push: image
#	docker push $(IMAGE_NAME)
