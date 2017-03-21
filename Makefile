SOURCES=appdependency.go

.PHONY: all templates darwin linux windows default

default: darwin

all: darwin linux windows

templates:
	zip -qr build-artifacts/templates.zip templates

darwin: $(SOURCES) templates
	GOOS=darwin go build -o build-artifacts/appdependency $(SOURCES)

linux: $(SOURCES) templates
	GOOS=linux go build -o build-artifacts/appdependency-linux $(SOURCES)

windows: $(SOURCES) templates
	GOOS=windows go build -o build-artifacts/appdependency.exe $(SOURCES)