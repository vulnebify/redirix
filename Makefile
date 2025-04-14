APP_NAME = redirix
CMD_PATH = ./cmd/redirix
BUILD_DIR = ./bin
VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: all build clean

all: build

build:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags "-X=github.com/vulnebify/redirix/internal/app.Version=$(VERSION)" \
	-o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

clean:
	rm -rf $(BUILD_DIR)
