.PHONY: build clean win linux mac
#.ONESHELL:

# go build -ldflags "-X main.Version=$(git describe --tags)"

build: clean win linux mac mac2
clean:
	rm -rf dist
	mkdir -p dist
win:
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -X main.Version=$$(git describe --tags)" -o ./dist/dev-server.exe main.go
	cd dist && 7z a -sdel dev-server-win32-x64.zip dev-server.exe
linux:
	GOOS=linux   GOARCH=amd64 go build -ldflags "-w -X main.Version=$$(git describe --tags)" -o ./dist/dev-server main.go
	cd dist && 7z a -sdel dev-server-linux-x64.tar dev-server
	cd dist && 7z a -sdel dev-server-linux-x64.tar.gz dev-server-linux-x64.tar
mac:
	GOOS=darwin  GOARCH=amd64 go build -ldflags "-w -X main.Version=$$(git describe --tags)" -o ./dist/dev-server main.go
	cd dist && 7z a -sdel dev-server-darwin-x64.tar dev-server
	cd dist && 7z a -sdel dev-server-darwin-x64.tar.gz dev-server-darwin-x64.tar
mac2:
	GOOS=darwin  GOARCH=arm64 go build -ldflags "-w -X main.Version=$$(git describe --tags)" -o ./dist/dev-server main.go
	cd dist && 7z a -sdel dev-server-darwin-arm64.tar dev-server
	cd dist && 7z a -sdel dev-server-darwin-arm64.tar.gz dev-server-darwin-arm64.tar
