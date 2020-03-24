FROM golang:1.14 as build
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOARCH=amd64
COPY . /app
RUN cd /app && go run github.com/gobuffalo/packr/v2/packr2 -v
RUN echo "package web \nimport _ \"github.com/AOEpeople/vistecture/v2/packrd\"" > /app/controller/web/web-packr.go
RUN cd /app && go build -o vistecture .
RUN ls -l /app

FROM alpine:3.9.5
RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini
COPY --from=build /app/vistecture /usr/local/bin
