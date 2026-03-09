CMD_NAME := check_akcp
SOURCES := $(shell find . -type f -name '*.go')

$(CMD_NAME): $(SOURCES)
	go build -trimpath -ldflags "-s -w" -o $(CMD_NAME) main.go

clean:
	go clean
	rm -f $(CMD_NAME)
