FROM golang:1.25-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o app ./cmd/api

CMD ["./app"]