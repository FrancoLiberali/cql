# builder image
FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN apk add build-base
RUN CGO_ENABLED=1 go build -a -o badaas .


# final image for end users
FROM alpine:3.16.2
COPY --from=builder /app/badaas .
EXPOSE 8000
ENTRYPOINT [ "./badaas" ]