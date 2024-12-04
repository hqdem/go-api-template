CFG_PATH=./configs/dev/common.yaml

.PHONY: build
build:
	go build -o app main.go

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
