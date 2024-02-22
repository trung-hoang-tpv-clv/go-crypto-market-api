# Go Crypto Market API

This project is a Go-based API for fetching cryptocurrency market prices.

## Project Structure

- `cmd`: Contains the main application entry points.
- `internal`: Contains the core logic of the application, structured into:
  - `config`: Configuration files and setup utilities for the database and Redis.
  - `domain`: Domain models used across the application.
  - `infrastructure`: External services and persistence layer implementations, including:
    - `cache`: Redis cache implementation.
    - `coingecko`: CoinGecko API client.
    - `repository`: Database repositories for storing and retrieving domain entities.
  - `interfaces`: Defines the application interfaces, including:
    - `dto`: Data Transfer Objects for external communication.
    - `handlers`: HTTP handlers for routing and handling requests.
    - `routers`: Setup of routing for the HTTP server.
  - `usecases`: Application business logic for handling requests and serving data.
- `pkg`: Additional packages that can be imported by external applications.

## Getting Started

To get the project up and running on your local machine for development and testing purposes, follow these steps:

### Prerequisites

- Go (version 1.15 or higher)
- Redis server running on the default port
- Access to a PostgreSQL database

Setting up Redis & PostgreSQL
- docker-compose up -d

### Installing

Clone the repository to your local machine:

```sh
git clone https://github.com/yourusername/go-crypto-market-api.git
cd go-crypto-market-api
go mod tidy

```
### Setup environment
```sh
cp .env.example .emv
after copy file .env.example to .env you can update your environment
```
### Running the application

```sh
go run cmd/main.go

