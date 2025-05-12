FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o card-service ./cmd/main.go
RUN ls -l /app
# Копируем entrypoint.sh и даём права
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 50052

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./card-service"]
