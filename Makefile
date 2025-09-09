build:
	GOOS=linux GOARCH=amd64 go build -o checker-website-links main.go

build-mac-darwin:
	GOOS=darwin GOARCH=amd64 go build -o checker-website-links main.go

build-mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -o checker-website-links main.go

build-windows-64:
	GOOS=windows GOARCH=amd64 go build -o checker-website-links main.go

build-windows-32:
	GOOS=windows GOARCH=386 go build -o checker-website-links main.go
	