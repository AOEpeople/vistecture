FROM alpine:3.5

RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini

RUN apk add --no-cache go

RUN apk add --no-cache git

RUN go get github.com/AOEpeople/vistecture

ENTRYPOINT ["vistecture"]