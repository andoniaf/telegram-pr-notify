.PHONY: build test test-v lint fmt vet fmt-check docker clean

BINARY := telegram-pr-notify
IMAGE  := telegram-pr-notify

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) .

test:
	go test ./...

test-v:
	go test -v ./...

lint: vet fmt-check

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:" && gofmt -l . && exit 1)

docker:
	docker build -t $(IMAGE) .

clean:
	rm -f $(BINARY)
