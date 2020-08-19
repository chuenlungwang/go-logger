GOCMD = go
GOBUILD = $(GOCMD) build
BINARY_NAME = bin
NAME = go-logger
VERSION = latest
RM = rm

.PHONY: all clean

all: build-linux build-macos build-windows

build-macos:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)/$(NAME)-mac-$(VERSION)

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)/$(NAME)-linux-$(VERSION)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)/$(NAME)-windows-$(VERSION)

clean:
	$(RM) -rf $(BINARY_NAME)