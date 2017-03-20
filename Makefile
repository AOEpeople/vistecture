SOURCES=appdependency.go

.PHONY: all darwin linux windows default

default: darwin

all: darwin linux windows


darwin: $(SOURCES)
	GOOS=darwin go build -o appdependency $(SOURCES)

linux: $(SOURCES)
	GOOS=linux go build -o appdependency-linux $(SOURCES)

windows: $(SOURCES)
	GOOS=windows go build -o appdependency.exe $(SOURCES)