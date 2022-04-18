defaut:
	go build
darwin:
	env GOOS=darwin GOARCH=amd64 go build
install:
	go install github.com/ycollatin/gocol
ed:
	vim *.go exempla/data/lemmes.* exempla/data/modeles.la exempla/data/morphos.fr
