FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/services/data

ENV HOST=${HOST}
ENV PORT=${PORT}
ENV JWT_SECRET=${JWT_SECRET}

RUN go build -o data-service main.go

CMD ["./data-service"]
