SOURCES=vistecture.go
VERSION=0.5.3

.PHONY: all templates darwin linux windows default

default: darwin

all: darwin linux windows

templates:
	zip -qr build-artifacts/templates.zip templates

darwin: $(SOURCES) templates
	GOOS=darwin go build -o build-artifacts/vistecture $(SOURCES)

linux: $(SOURCES) templates
	GOOS=linux go build -o build-artifacts/vistecture-linux $(SOURCES)

windows: $(SOURCES) templates
	GOOS=windows go build -o build-artifacts/vistecture.exe $(SOURCES)

dockerpublish:
	docker build --no-cache -t aoepeople/vistecture .
	docker tag aoepeople/vistecture:latest aoepeople/vistecture:$(VERSION)
	docker push aoepeople/vistecture