APP      := rss-proxy
PORT     := 8000
CONFIG   := config.yml

.PHONY: help test run build clean docker-build docker-run

help:
	@echo "Targets disponibles:"
	@echo "  make test        → lance tous les tests"
	@echo "  make run         → lance le serveur en local"
	@echo "  make build       → build le binaire"
	@echo "  make clean       → supprime le binaire"
	@echo "  make docker-build→ build l'image Docker"
	@echo "  make docker-run  → run l'image Docker"

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

docker-run:
	docker run --rm -p $(PORT):8000 -v $(PWD)/$(CONFIG):/app/$(CONFIG) $(APP)

