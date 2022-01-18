all: build

build: linux mac win

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/vuitton cmd/main.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/vuitton_darwin cmd/main.go

win:
	GOOS=windows GOARCH=amd64 go build -o bin/vuitton.exe cmd/main.go

lint:
	staticcheck ./...
