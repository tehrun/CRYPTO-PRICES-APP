# Cryptocurrency Prices App

## Overview
The Cryptocurrency Prices App is a Go-based application that fetches and displays real-time cryptocurrency prices. It provides a simple API for users to retrieve price information for various cryptocurrencies.

## Features
- Fetch real-time cryptocurrency prices from external APIs.
- API endpoints to get prices for all cryptocurrencies or a specific one.
- Configurable settings for API keys and other parameters.
- Continuous Integration and Continuous Deployment pipelines for automated testing and deployment.

## Project Structure
```
crypto-prices-app
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── api
│   │   └── handlers.go      # HTTP handlers for API requests
│   ├── config
│   │   └── config.go        # Configuration management
│   ├── prices
│   │   └── service.go       # Business logic for price fetching
│   └── types
│       └── types.go         # Custom types and structures
├── pkg
│   └── client
│       └── http.go          # HTTP client for external API requests
├── .github
│   └── workflows
│       ├── ci.yml           # CI pipeline configuration
│       └── cd.yml           # CD pipeline configuration
├── go.mod                    # Module definition
├── go.sum                    # Dependency checksums
└── README.md                 # Project documentation
```

## Setup Instructions
1. Clone the repository:
   ```
   git clone https://github.com/yourusername/crypto-prices-app.git
   cd crypto-prices-app
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Configure your environment variables or create a configuration file as needed.

4. Run the application:
   ```
   go run cmd/server/main.go
   ```

## Usage
- Access the API at `http://localhost:8080/api/prices` to get the list of cryptocurrency prices.
- Use `http://localhost:8080/api/prices/{id}` to get the price of a specific cryptocurrency by its ID.

## CI/CD
This project includes CI/CD pipelines defined in the `.github/workflows` directory. The CI pipeline runs tests and builds the application on every push, while the CD pipeline deploys the application to production after successful builds.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.