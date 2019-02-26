FROM golang:1.11 as build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go get -u github.com/gobuffalo/packr/v2/packr2
RUN mkdir -p /go/src/github.com/AOEpeople/
COPY . /go/src/github.com/AOEpeople/vistecture
RUN cd /go/src/github.com/AOEpeople/vistecture && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go mod vendor
RUN cd /go/src/github.com/AOEpeople/vistecture && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr2 -v
RUN cd /go/src/github.com/AOEpeople/vistecture && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
RUN ls -al /go/src/github.com/AOEpeople/vistecture

FROM alpine:3.6
RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini
COPY --from=build /go/src/github.com/AOEpeople/vistecture/vistecture /usr/local/bin
