FROM golang:1.22.0-alpine

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o kafcmd ./cmd/kafcmd

CMD ["./kafcmd"]
