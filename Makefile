build:
	go build -o ./bin/gogh cmd/gogh/main.go

dev:
	find . -name "*.go" | entr -r make build
