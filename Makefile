version := 0.0.0
.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Makefile Commands:"
	@echo "----------------------------------------------------------------"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo "----------------------------------------------------------------"

run: ## run server
	@go run main.go

.PHONY: up
up: ## start docker container
	@docker-compose -f docker-compose.yml pull
	@docker-compose -f docker-compose.yml up -d

.PHONY: down
down: ## shut down docker container
	docker-compose -f docker-compose.yml down --remove-orphans

docker-build: ## build docker image
	docker build -t colemanword/thermomatic:$(version) .

docker-push: ## push docker image
	docker push colemanword/thermomatic:$(version)
	docker tag colemanword/thermomatic:$(version) colemanword/thermomatic:latest
	docker push colemanword/thermomatic:latest

test: ## run tests
	@go test -v ./...