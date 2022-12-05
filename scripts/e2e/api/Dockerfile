# builder image
FROM golang:1.19-alpine as builder
RUN apk add build-base
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 go build --race  -a -o badaas .

FROM alpine:3.16.2
ENV BADAAS_PORT=8000
COPY --from=builder /app/badaas .
COPY ./scripts/e2e/api/badaas.yml . 
ENTRYPOINT [ "./badaas" ]
