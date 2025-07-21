# CheckoutApi

This project is a robust, production-ready implementation of a supermarket checkout system API. It is designed with a clean, decoupled architecture to support scalability and frequent business rule changes, such as modifying item prices or special offers.

The core of the system is built around a stateful checkout session that strictly implements the `ICheckout` interface, while the surrounding architecture provides modern features like dynamic configuration and a RESTful API.

## Key Features

- **RESTful API**: Provides endpoints to create checkout sessions, scan items, and retrieve totals.
- **Complex Pricing Logic**: Natively handles individual item prices and multi-buy special offers (e.g., "3 for $130").
- **Dynamic Configuration**: Pricing rules are loaded from an external `pricing.json` file, completely decoupling business rules from compiled code.
- **Hot-Reloading**: The server automatically detects changes to `pricing.json` and updates its pricing rules **without requiring a restart**, demonstrating a high-availability design pattern.
- **Clean Architecture**: The code is organized into distinct layers (Domain, Repository, Handler) to ensure separation of concerns, high cohesion, and low coupling.
- **Comprehensive Test Suite**: Includes unit tests for core logic and integration tests for HTTP handlers, ensuring code quality and reliability.
- **Concurrency Safe**: The application is designed to handle multiple concurrent requests safely using mutexes for shared resources.

## Architecture Overview

The project follows the principles of **Clean Architecture**. The core business logic is isolated in the `checkout` and `domain` packages, with no knowledge of the outside world (like HTTP or databases).

```
checkoutapi/
├── checkout/
│   ├── handler/
│   │   ├── http.go
│   │   └── http_test.go
│   ├── repository/
│   │   ├── memory.go
│   │   └── memory_test.go
│   ├── checkout.go
│   └── checkout_test.go
├── cmd/
│   └── checkoutapi/
│       ├── configs/
│       │   └── pricing.json
│       └── checkoutapi.go
├── domain/
│   └── checkout.go
├── pricing/
│   └── service/
│       └── service.go
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- Go (version 1.22 or later)

### Installation & Setup

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/TheFodfather/checkoutapi
    cd checkout-system
    ```

2.  **Create the configuration file:**
    Create a file at `configs/pricing.json` with the following content. The server will read from this file at startup.

    ```json
    {
      "A": {
        "unitPrice": 50,
        "specialPrice": {
          "quantity": 3,
          "price": 130
        }
      },
      "B": {
        "unitPrice": 30,
        "specialPrice": {
          "quantity": 2,
          "price": 45
        }
      },
      "C": {
        "unitPrice": 20,
        "specialPrice": null
      },
      "D": {
        "unitPrice": 15,
        "specialPrice": null
      }
    }
    ```

3.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

### Running the Application

To run the API server, execute the following command from the root of the project:

```sh
go run ./cmd/checkoutapi/checkoutapi.go
```

You should see output indicating that the server has started and loaded the pricing rules:

```
✅ Successfully loaded new pricing rules.
🚀 Starting server on http://localhost:8080
```

---

## API Usage

Full API documentation is available here: **[api.md](api.md)**

---

## Hot-Reloading Demonstration

This feature allows you to change prices on the fly.

1.  **Start the server** and create a new checkout session as shown above.

2.  **Scan an item** and check its price.

3.  **Modify the configuration**. While the server is running, open `configs/pricing.json` and change the price of item `A`.
    Save the file. You will see a log message in the server console: `🔄 Change detected...`.

4.  **Scan the same item again and then get the total price**.
    The new total will be the original price + the update price, without the need to restart the server.

## Running the Tests

The project has a comprehensive test suite.

- **Run all tests** for all packages:

  ```sh
  go test ./... -v
  ```

- **Run tests for a specific package**:

  ```sh
   go test ./checkout -v
  ```

- **Run tests with the race detector** to check for concurrency issues. This requires `cgo` to be enabled.
  ```sh
  CGO_ENABLED=1 go test ./... -race
  ```
