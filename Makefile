VERSION = 1.0
TAG = $(VERSION)
PREFIX = tommymuehle/speakerdeck_exporter
GIT_COMMIT = $(shell git rev-parse --short HEAD)

container:
	docker build --build-arg VERSION=$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(PREFIX):$(TAG) .

lint:
	golangci-lint run
