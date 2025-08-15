APP_NAME = chrono-ntp
SRC = src/main.go

.PHONY: all test build clean run readme-demo

all: test build

test:
	go test ./src/...

build:
	go build -o $(APP_NAME) $(SRC)

run:
	go run $(SRC) --server=time.google.com

clean:
	rm -f $(APP_NAME)

readme-demo:
	vhs README-demo.tape --output README-demo.gif
