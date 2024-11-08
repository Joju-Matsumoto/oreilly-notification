build: ## go build
	go build -o oreilly-notification ./cmd/oreilly-notification

test: ## go test
	go test -v ./...

generate: ## go generate
	go generate ./...