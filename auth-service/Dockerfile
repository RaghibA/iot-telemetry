FROM golang:1.23-alpine

WORKDIR /app

ARG PORT

ENV HOST=${HOST}
ENV PORT=${PORT}
ENV JWT_SECRET=${JWT_SECRET}

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/server

CMD ["./auth-service"]
