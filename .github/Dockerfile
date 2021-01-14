FROM alpine:latest
RUN apk add --no-cache \
  graphviz \
  font-bitstream-type1 \
  inotify-tools \
  tini
COPY build-artifacts/vistecture-linux /usr/local/bin/vistecture
