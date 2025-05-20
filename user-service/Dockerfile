FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o user-service ./cmd/user-service

# Копируем entrypoint.sh и даём права
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 50051

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./user-service"]
