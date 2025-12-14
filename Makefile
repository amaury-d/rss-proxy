APP      := rss-proxy
PORT     := 8000
CONFIG   := config.yml

.PHONY: help test run build clean docker-build docker-run

help:
	@echo "Targets disponibles:"
	@echo "  make test        → run all tests"
	@echo "  make run         → run the server locally"
	@echo "  make build       → build binary"
	@echo "  make clean       → clean project"
	@echo "  make docker-build→ build docker image"
	@echo "  make docker-run  → run docker image"

test:
	go test ./...

run:
	go run main.go

build:
	go build -o $(APP)

clean:
	rm -f $(APP)

docker-build:
	docker build -t $(APP) .

docker-run: docker-build
	docker run --rm -p $(PORT):8000 -v $(PWD)/$(CONFIG):/app/$(CONFIG) $(APP)

