FROM golang:1.20.3-alpine as dev

RUN go install github.com/cespare/reflex@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /code

FROM dev as builder

COPY . /code

RUN go build -o ./run .


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /code/run .

EXPOSE 8080
CMD ["./run"]
