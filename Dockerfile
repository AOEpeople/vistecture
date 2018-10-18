FROM golang:1.9 as build
RUN go get github.com/AOEpeople/vistecture

FROM alpine:latest
RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini
COPY --from=build /go/bin/vistecture /usr/local/bin
