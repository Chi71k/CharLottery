# Card Service

This project implements a gRPC-based card service using Go and MongoDB. It provides functionalities for managing card data, including creating, retrieving, and validating cards.

## Project Structure

```
card-service
├── cmd
│   └── main.go                # Entry point of the application
├── pkg
│   ├── api
│   │   └── card_service.proto  # gRPC protocol definitions
│   ├── db
│   │   ├── mongo.go            # MongoDB initialization
│   │   └── models
│   │       └── card.go         # Card model definition
│   ├── service
│   │   └── card_service.go      # gRPC service logic
├── go.mod                       # Go module definition
├── go.sum                       # Module dependency checksums
└── README.md                    # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- MongoDB

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd card-service
   ```

2. Install the necessary dependencies:
   ```
   go mod tidy
   ```

### Running the Service

1. Start your MongoDB server.

2. Run the application:
   ```
   go run cmd/main.go
   ```

### API Documentation

Refer to the `card_service.proto` file in the `pkg/api` directory for the gRPC service definitions and message structures.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or features.

### License

This project is licensed under the MIT License. See the LICENSE file for details.