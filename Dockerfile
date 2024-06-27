FROM node:16.20.2-alpine3.18 AS frontend
COPY . /app
RUN apk add --update python3 make gcc g++
RUN cd /app/controller/web/template && npm install && npm run build

FROM golang:1.17 AS build
ENV GOOS=linux
ENV CGO_ENABLED=0
COPY . /app
COPY --from=frontend /app/controller/web/template/dist /app/controller/web/template/dist/
RUN cd /app && go build -o vistecture .

FROM alpine:latest
RUN apk add --no-cache \
  graphviz \
  ttf-freefont \
  inotify-tools \
  tini
COPY --from=build /app/vistecture /usr/local/bin
