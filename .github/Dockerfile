FROM alpine:latest
RUN apk add --no-cache \
  graphviz \
  ttf-freefont \
  inotify-tools \
  tini
COPY build-artifacts/vistecture-linux /usr/local/bin/vistecture
