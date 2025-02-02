# Purchase cart service

The Purchase Cart Service is a RESTful API designed to handle purchase cart operations. It accepts requests with order details, calculates pricing information (including VAT), and returns a structured response with details for the order and individual items.

## How to run

1. Ensure you have **Docker** installed on your system.
2. Clone the Git repository for this project:
   ```bash
   git clone https://github.com/ImValerio/purchase-cart-service.git
   cd purchase-cart-service
   ```
3. For a fast startup:\
    `docker compose up`\
    \
    You can also build the image\
   `docker build -t purchase-cart-service`\
   and run the specific script:\
    `docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt purchase-cart-service ./scripts/run.sh`

## Considerations

### Technology Stack

Using **Golang** as the backend programming language, ensuring high performance and efficient handling of concurrent operations.\
For persistent storage, the service uses **Postgres**, a robust and reliable relational database system that facilitates efficient management of pricing data, orders, and related information.

### Structure

A modular approach has been followed for better code readability and maintainability. The application is structured into distinct layers, including handlers, services, and repositories, each with clearly defined responsibilities.

### Testing

Integration tests are available using the command (db must be running):\
`docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt purchase-cart-service ./scripts/test.sh`

### Potential improvements

- Implement an in-memory caching solution to reduce database load and improve response times for frequently accessed product data.
- Ensure atomicity by wrapping critical database operations in transactions to prevent partial updates in case of failures.
- Implement a circuit breaker to handle failures gracefully and avoid cascading failures by limiting requests to failing services.
- Rate Limiting 
