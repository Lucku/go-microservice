FROM golang:1.12 as builder
COPY . /app
WORKDIR /app
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o otto-rest-api ./cmd/toysapiserver/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app .
CMD ["./otto-rest-api"]