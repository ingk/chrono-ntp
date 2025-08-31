APP_NAME = chrono-ntp

.PHONY: all test build clean run readme-demo

all: test build

test:
	go test ./src/...

build:
	go build -o $(APP_NAME) ./src/...

build-raspberry-arm64:
	mkdir -p dist && rm -f dist/$(APP_NAME)-raspberry-arm64 && go build -o dist/$(APP_NAME)-raspberry-arm64 ./src/...

run:
	go run ./src/... --server=time.google.com

clean:
	rm -f $(APP_NAME)

readme-demo:
	vhs assets/demo.tape --output assets/demo.gif
