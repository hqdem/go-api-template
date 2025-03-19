CFG_PATH=./configs/dev/common.yaml

.PHONY: build
build:
	go build -o app main.go

.PHONY: swagger
swagger:
	swag init -g internal/commands/runserver/runserver.go && swag fmt

.PHONY: runserver
runserver:
	./app run --config $(CFG_PATH)

.PHONY: staticcheck
staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: bin-deps
bin-deps:
	go install github.com/swaggo/swag/cmd/swag@latest