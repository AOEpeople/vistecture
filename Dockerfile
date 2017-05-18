FROM alpine:3.5

RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini

RUN apk add --no-cache gcc && apk add --no-cache libc-dev

RUN apk add --no-cache go

RUN apk add --no-cache git

RUN mkdir -p /usr/src/go

RUN export CGO_ENABLED=0 && export GOPATH="/usr/src/go" && go get github.com/AOEpeople/vistecture  && cp -R /usr/src/go/src/github.com/AOEpeople/vistecture/templates /usr/src/go/bin

ENV PATH "$PATH:/usr/src/go/bin"