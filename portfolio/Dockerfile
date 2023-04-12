FROM golang:1.20.3-alpine as builder

WORKDIR /code

# Get Reflex for live reload in dev
RUN go install github.com/cespare/reflex@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./run .


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

#Copy executable from builder
COPY --from=builder /code/run .

EXPOSE 8080
CMD ["./run"]
