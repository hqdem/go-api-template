CFG_PATH=./configs/dev/common.yaml

build:
	go build -o app main.go

runserver:
	./app run --config $(CFG_PATH)

staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

lint:
	golangci-lint run ./...
