FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/services/admin

RUN go build -o admin-service main.go

CMD ["./admin-service"]
