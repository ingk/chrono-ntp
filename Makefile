APP_NAME = chrono-ntp
SRC = main.go

.PHONY: all build clean run

all: build

build:
	go build -o $(APP_NAME) $(SRC)

run:
	go run $(SRC) --server=time.google.com

clean:
	rm -f $(APP_NAME)
