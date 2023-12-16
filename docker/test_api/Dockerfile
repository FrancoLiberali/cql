# DEVELOPMENT ONLY, DO NOT USE FOR PRODUCTION
FROM golang:1.19-alpine
RUN addgroup -S badaas \
    && adduser -S badaas -G badaas
USER badaas
WORKDIR /badaas
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS=-buildvcs=false