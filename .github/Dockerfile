FROM alpine:latest
ARG TARGETARCH
RUN apk add --no-cache \
  graphviz \
  ttf-freefont \
  inotify-tools \
  tini
COPY --chmod=0755 build-artifacts/vistecture-linux-$TARGETARCH /usr/local/bin/vistecture
