FROM golang:1.18.1-alpine3.15


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
