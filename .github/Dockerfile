FROM alpine:latest
ARG TARGETARCH
RUN apk add --no-cache \
  graphviz \
  ttf-freefont \
  inotify-tools \
  tini
COPY build-artifacts/vistecture-linux-$TARGETARCH /usr/local/bin/vistecture
