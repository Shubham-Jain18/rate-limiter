# Rate Limiter

This project demonstrates four rate-limiting algorithms using Redis and Gin in Go.

## Implemented Algorithms
- Sliding Window Counter
- Token Bucket
- Leaky Bucket
- Sliding Window Log

## How to Run

1. **Start Redis (using Docker)**:  
   ```
   docker-compose up -d
   ```  
2. **Run the server (default port is 8080)**:
   ```  
   go run main.go
   ``` 
   Or use a different port to simulate another instance:  
   ```
   go run main.go 8081
   ```  

## Endpoints to Test
You can hit the following endpoints:
```
- GET /ping/sliding-window-counter
- GET /ping/token-bucket
- GET /ping/leaky-bucket
- GET /ping/sliding-window-log
```

Example:  
```
curl http://localhost:8080/ping/leaky-bucket
```

## Simulating Distributed Systems
To simulate multiple services (like in distributed systems), launch multiple server instances on different ports:  
  ```
  go run main.go 8081  
  go run main.go 8082
  ``` 
Since all instances connect to the same Redis backend, rate-limiting is shared across instances.  
