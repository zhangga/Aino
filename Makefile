# 当前目录
CUR_DIR=$(shell pwd)
OUT_DIR=$(CUR_DIR)/bin

# 命令
GO_BUILD = CGO_ENABLED=0 go build -trimpath
GO_RUN = CGO_ENABLED=0 go run

SERVER_VERSION	?= $(shell git describe --long --tags --dirty --always)
SERVER_VERSION	?= unkonwn
BUILD_TIME      ?= $(shell date "+%F_%T_%Z")
COMMIT_SHA1     ?= $(shell git show -s --format=%h)
COMMIT_LOG      ?= $(shell git show -s --format=%s)
COMMIT_AUTHOR	?= $(shell git show -s --format=%an)
COMMIT_DATE		?= $(shell git show -s --format=%ad)
COMMIT_MSG		?= $(COMMIT_AUTHOR)|$(COMMIT_DATE)|${COMMIT_LOG}
VERSION_PACKAGE	?=github.com/zhangga/aino/pkg/version

CUR_BRANCH := $(shell git branch --show-current)


VERSION_BUILD_LDFLAGS= \
-X "${VERSION_PACKAGE}.Version=${SERVER_VERSION}" \
-X "${VERSION_PACKAGE}.BuildTime=${BUILD_TIME}" \
-X "${VERSION_PACKAGE}.CommitHash=${COMMIT_SHA1}" \
-X "${VERSION_PACKAGE}.Description=${COMMIT_MSG}"
.PHONY: build
# build
build:
	$(GO_BUILD) \
	-ldflags '$(VERSION_BUILD_LDFLAGS)' \
	-o $(OUT_DIR)/ \
	./cmd/aino

.PHONY: run
# run
run:
	$(GO_RUN) \
	-ldflags '$(VERSION_BUILD_LDFLAGS)' \
	./cmd/aino run

.PHONY: test
# run all test
test:
	@echo "Running tests..."
	@go test -v -race ./...

.PHONY: lint
# run all lint
lint:
	golangci-lint run -c .golangci.yml ./...

.PHONY: milvus
# run milvus
milvus:
	docker compose -f docker-compose/milvus-standalone-docker-compose.yml up -d

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help