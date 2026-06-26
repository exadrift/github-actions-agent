NAME := agent
GO_VERSION := 1.26
LINTER_VERSION := 2.12.2
VERSION ?= local
MAJOR_VERSION ?= $(shell echo $(VERSION) | cut -d . -f1)
GO_IMAGE := golang:$(GO_VERSION)-alpine
GO_RUN := docker run --rm -e CGO_ENABLED=0 -e HOME=$$HOME -v $$HOME:$$HOME -u $(shell id -u):$(shell id -g) -v $(shell pwd):/build -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock -w /build $(GO_IMAGE) go
GO_FILES := $(shell find . -type f -path **/*.go -not -path "./vendor/*")
IMAGE_NAME := exadrift/$(NAME)
IMAGE_TAG := $(IMAGE_NAME):$(VERSION)

.PHONY: test
test:
	$(GO_RUN) test -cover -p 1 --timeout 10m ./...

.PHONY: lint-check
lint-check:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v$(LINTER_VERSION) golangci-lint run

.PHONY: build
build: build-docker

.PHONY: build-docker
build-docker: bin/$(NAME)
	docker build --build-arg NAME=$(NAME) -t $(IMAGE_TAG) .

.PHONY: publish
publish: build-docker
	docker push $(IMAGE_TAG)

bin/$(NAME): $(GO_FILES)
	$(GO_RUN) build -trimpath -ldflags="-s -w -X 'main.Version=$(VERSION)'" -mod=vendor -o ./bin/$(NAME) *.go

.PHONY: install
install: bin/$(NAME)
	cp ./bin/$(NAME) /usr/local/bin/

.PHONY: clean
clean:
	rm -rf bin
	docker image rm -f $(IMAGE_TAG)
