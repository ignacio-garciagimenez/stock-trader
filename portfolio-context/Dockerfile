FROM golang:1.20.3-alpine as builder

WORKDIR /code

COPY . .

RUN go build -o ./run .



FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /code/run .

EXPOSE 8080
CMD ["./run"]
