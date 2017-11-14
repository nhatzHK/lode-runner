GOCMD   := go
GOBUILD := $(GOCMD) build
GOTEST  := $(GOCMD) test

.PHONY: build test certs

build:
	$(GOBUILD) -o build/server
	sudo setcap cap_net_bind_service=+ep build/server
	cp -r game/levels/ build/

test:
	$(GOTEST) -coverprofile=coverage.out

certs:
	openssl req -new -nodes -x509 -newkey rsa -keyout build/server.key -out build/server.crt -subj "/"
