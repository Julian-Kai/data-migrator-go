## Run lint in current directory and all of its subdirectories
.PHONY: build
build:
	@go build -o app cmd/server/main.go
