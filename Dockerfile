FROM node:15.5.1-alpine3.12 as frontend
COPY . /app
RUN cd /app/controller/web/template && npm install && npm run build

FROM golang:1.15 as build
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOARCH=amd64
COPY . /app
COPY --from=frontend /app/controller/web/template/dist /app/controller/web/template/dist/
RUN cd /app && go run github.com/gobuffalo/packr/v2/packr2 -v
RUN echo "package web \nimport _ \"github.com/AOEpeople/vistecture/v2/packrd\"" > /app/controller/web/web-packr.go
RUN cd /app && go build -o vistecture .
RUN ls -l /app

FROM alpine:latest
RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini
COPY --from=build /app/vistecture /usr/local/bin
