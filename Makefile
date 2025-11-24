APP_NAME = chrono-ntp

.PHONY: all test build clean run readme-demo

all: test build

test:
	go test ./...

build:
	go build -o $(APP_NAME) main.go

run:
	go run ./... --server=time.google.com

clean:
	rm -f $(APP_NAME)

readme-demo:
	vhs assets/demo.tape --output assets/demo.gif
