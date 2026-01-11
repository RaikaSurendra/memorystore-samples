# Go Caching Sample

This directory contains the Go implementation of the caching sample.

## Directory Structure

-   `sample-demo-app/`: A full Go Caching application using Gin, Redis (go-redis), and PostgreSQL (pgx).

## How to Run

1.  Navigate to `sample-demo-app`.
2.  Ensure you have a Redis and PostgreSQL instance running.
3.  Set environment variables:
    ```bash
    export REDIS_HOST=localhost
    export DB_HOST=localhost
    export DB_PASS=yourpassword
    ```
4.  Run the application:
    ```bash
    go run main.go
    ```
5.  Test endpoints:
    -   `POST /items`: Create an item.
    -   `GET /items/:name`: Retrieve an item (check source in response to see if from cache or database).
