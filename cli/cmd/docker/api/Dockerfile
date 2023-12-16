# builder image
FROM golang:1.19-alpine AS builder
RUN apk add build-base
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build -o badaas_service -a .


# final image for end users
FROM alpine:3.16.2
RUN addgroup -S badaas \
    && adduser -S badaas -G badaas
USER badaas
COPY --from=builder /app/badaas_service .
COPY ./badaas/config/badaas.yml .
EXPOSE 8000
ENTRYPOINT ["./badaas_service", "--config_path", "badaas.yml"]