CFG_PATH=./configs/dev/common.yaml

build:
	go build -o app main.go

runserver:
	./app run --config $(CFG_PATH)
