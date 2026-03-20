BINARY_NAME=mailtrap
VERSION ?= dev

.PHONY: build install clean

build:
	go build -ldflags "-s -w" -o $(BINARY_NAME) .

install:
	go install .

clean:
	rm -f $(BINARY_NAME)
