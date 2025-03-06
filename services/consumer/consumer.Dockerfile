FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/services/consumer

ENV HOST=${HOST}
ENV PORT=${PORT}
ENV JWT_SECRET=${JWT_SECRET}

RUN go build -o consumer-service main.go

CMD ["./consumer-service"]
