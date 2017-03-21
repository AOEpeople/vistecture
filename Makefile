SOURCES=appdependency.go

.PHONY: all darwin linux windows default

default: darwin

all: darwin linux windows


darwin: $(SOURCES)
	GOOS=darwin go build -o build-artifacts/appdependency $(SOURCES)

linux: $(SOURCES)
	GOOS=linux go build -o build-artifacts/appdependency-linux $(SOURCES)

windows: $(SOURCES)
	GOOS=windows go build -o build-artifacts/appdependency.exe $(SOURCES)