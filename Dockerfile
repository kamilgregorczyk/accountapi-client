FROM golang:1.15.5-alpine

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD ["./entrypoint.sh"]
