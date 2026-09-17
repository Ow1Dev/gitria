.PHONY: run build clean

APP_DIR := ./cmd/gitria

export SSH_LISTEN_PORT := 2222
export SSH_HOST_KEY := ./id_rsa

run:
	go run $(APP_DIR)

build:
	mkdir -p bin
	go build -o bin/gitria $(APP_DIR)

clean:
	rm -rf bin

