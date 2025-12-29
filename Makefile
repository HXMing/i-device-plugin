IMG?=docker.io/hongxuming/i-device-plugin:latest

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/i-device-plugin cmd/main.go

.PHONY: build-image
build-image:
	docker build -t ${IMG} .

.PHONY: clean
clean:
	rm -rf bin/