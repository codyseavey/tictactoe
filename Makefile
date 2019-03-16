DOCKER_REPO=belligerence/tictactoe
APP_NAME=tictactoe
GIT_REPO=github.com/codyseavey/tictactoe

.PHONY: help
help: ## - Displays help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: ls
ls: ## - List images that were created
	@docker image ls $(DOCKER_REPO)

.PHONY: pullcerts
pullcerts: ## - Pulls certs from credstash
	@mkdir -p certs
	@credstash get cert.pem > certs/cert.pem
	@credstash get chain.pem > certs/chain.pem
	@credstash get fullchain.pem > certs/fullchain.pem
	@credstash get privkey.pem > certs/privkey.pem

.PHONY: build
build: ## - Docker build
	@docker build -t $(DOCKER_REPO) .

.PHONY: build-no-cache
build-no-cache: ## - Docker build with no-cache setting
	@docker build --no-cache -t $(DOCKER_REPO) .

.PHONY: run
run: stop pullcerts ## - Run the container that was built
	@docker run --name $(APP_NAME)-db -e POSTGRES_USER=$(APP_NAME) -e POSTGRES_PASSWORD=$(APP_NAME) -e POSTGRES_DB=$(APP_NAME) -d postgres
	@docker run --name $(APP_NAME) -d -p 443:8443 --link $(APP_NAME)-db:postgres -v $(PWD)/certs:/certs $(DOCKER_REPO):latest

.PHONY: test
test: pullcerts ## - Run tests for the application that was built
	@docker run --name $(APP_NAME)-test-db -e POSTGRES_USER=$(APP_NAME) -e POSTGRES_PASSWORD=$(APP_NAME) -e POSTGRES_DB=$(APP_NAME) -d postgres &> /dev/null
	@docker run --rm --link $(APP_NAME)-test-db:postgres -v $(GOPATH):/go -w /go/src/$(GIT_REPO) -e GOARCH=386 golang:alpine go test
	@docker rm -f $(APP_NAME)-test-db &> /dev/null || true

.PHONY: stop
stop: ## - Removes the container if it is running
	@docker rm -f $(APP_NAME) $(APP_NAME)-db 2> /dev/null || true

.PHONY: publish
publish: ## - Pushes the image to docker registry
	@docker push $(DOCKER_REPO):latest
