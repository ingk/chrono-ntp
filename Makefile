APP_NAME = chrono-ntp
SRC = src/main.go

.PHONY: all build clean run readme-demo

all: build

build:
	go build -o $(APP_NAME) $(SRC)

run:
	go run $(SRC) --server=time.google.com

clean:
	rm -f $(APP_NAME)

readme-demo:
	vhs README-demo.tape --output README-demo.gif
