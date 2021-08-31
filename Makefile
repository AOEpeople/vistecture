SOURCES=vistecture.go
VERSION=2.5.0

.PHONY: all templates darwin linux frontend windows default

default: darwin

demo: frontend run-example

run-example:
	go run vistecture.go --config example/demoproject/project.yml serve

all: frontend darwin_binary linux_binary windows_binary

windows: frontend windows_binary

linux: frontend linux_binary

darwin: frontend darwin_binary

frontend:
	cd ./controller/web/template && npm install && npm run build

templates:
	mkdir -p build-artifacts
	zip -qr build-artifacts/templates.zip templates

darwin_binary: $(SOURCES) templates
	GOOS=darwin go build -o build-artifacts/vistecture $(SOURCES)

linux_binary: $(SOURCES) templates
	GOOS=linux go build -o build-artifacts/vistecture-linux $(SOURCES)

windows_binary: $(SOURCES) templates
	GOOS=windows go build -o build-artifacts/vistecture.exe $(SOURCES)

docker:
	docker build --no-cache -t aoepeople/vistecture .
	docker tag aoepeople/vistecture:latest aoepeople/vistecture:$(VERSION)

docker-publish:
	docker push aoepeople/vistecture:latest
	docker push aoepeople/vistecture:$(VERSION)

dockerpublishexampleproject:
	cd example && ./generate-docs-with-docker.sh
	cd example && docker build --no-cache -t aoepeople/vistecture-example .
	docker push aoepeople/vistecture-example:latest
