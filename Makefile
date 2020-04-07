defaut:
	go build
darwin:
	env GOOS=darwin GOARCH=amd64 go build
install:
	go install github.com/ycollatin/gocol
