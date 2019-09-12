SOURCES=vistecture.go
VERSION=2.0.9

.PHONY: all templates darwin linux windows default

default: darwin

all: darwin linux windows

templates:
	packr2
	echo "package web \n import _ \"github.com/AOEpeople/vistecture/v2/packrd\"" > controller/web/web-packr.go
	mkdir -p build-artifacts
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
	docker push aoepeople/vistecture:latest
	docker push aoepeople/vistecture:$(VERSION)



dockerpublishexampleproject:
	cd example && ./generate-docs-with-docker.sh
	cd example && docker build --no-cache -t aoepeople/vistecture-example .
    docker push aoepeople/vistecture-example:latest