#!/usr/bin/zsh

export SERVERS="localhost:19101" \
  && export REDIS_USERNAME="" \
  && export REDIS_PASSWORD="123456" \
  && make build \
  && ./cmd/main