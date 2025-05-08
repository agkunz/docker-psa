.PHONY: build install clean

build:
	go build -o docker-psa ./cmd/

install: build
	mkdir -p ~/.docker/cli-plugins
	cp docker-psa ~/.docker/cli-plugins/docker-psa
	chmod +x ~/.docker/cli-plugins/docker-psa
	@echo "Docker PSA plugin installed. You can now run 'docker psa'"

clean:
	rm -f docker-psa
	rm -f ~/.docker/cli-plugins/docker-psa

test:
	go test ./...

.DEFAULT_GOAL := build