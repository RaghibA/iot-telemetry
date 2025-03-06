FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/cmd/migration

RUN go build -o migration main.go

CMD ["./migration"]