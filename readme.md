# Canonical Audit Log Service 
This is a mock API intended to capture and query audit events. Built with Golang and MongoDB, it leverages JWT for authentication.

## Overview
The API allows users to:

- Obtain a **JWT** token for protected routes.
- **Create** an audit log event.
- **Query** audit log events based on specific parameters.

## Getting started

### Installation 
This is a **dockerized** service. 

1. Clone repository 
2. Cd to **./canonicalAuditlog**
3. run `./deploy.sh` in terminal
4. The shell script is designed to set up the service in **ubuntu**.
5. The script installs **docker** and **docker-compose** if not already installed. 
6. A **docker daemond** is also spun up with `sudo systemctl start docker`
7. A **binary** is also built and run automatically.
8. **Dockerfile** handles the remaining dependencies.

## API

### JWT - GET
The API uses **JWT authentication**. A token must be acquired first by hitting the **/jwt** endpoint. \
This token must then be used in the Authorization header with the **"Bearer"** prefix for other API calls. 

### Usage
```bash
curl -X GET http://localhost:8080/jwt
```
```bash
# Example Response
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg2ODcyODl9.xUgLlNSCMOJsxFJwWDvl8RiMhMaY73tFnkezp7bGacw"}
```
This **JWT token** will be used to authenticate the remaining endpoints. 

### EVENTS/CREATE - POST
Creates a new audit log event. Requires a **JWT token** in the Authorization header.

### Usage
```bash
curl -X POST http://localhost:8080/events/create \
-H "Authorization: Bearer <TOKEN>" \
-H "Content-Type: application/json" \
-d '{"event_type": "type", "service_name": "service", "customer_id": "12345", "event_details": {"key": "value"}}'
```
```bash 
curl -X POST http://localhost:8080/events/create \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc0MjQwODJ9.yRJklHq-UqUENmgdpB4UEopfT9JiJKL9L3T-fS4VCZA" \
     -H "Content-Type: application/json" \
     -d '{
          "event_type": "item_review",
          "service_name": "review_service",
          "customer_id": "cus_118",
          "event_details": {
              "action": "submitted_review",
              "item_id": "item892",
              "rating": 0,
              "comment": "Bad"
          }
     }' 
```
```bash
# Example Response
200
```
**All** fields shown in the request are required. 

### EVENTS/QUERY - GET
Retrieves events based on query parameters. Requires a **JWT token** in the Authorization header.

### Usage
```bash
curl -X GET http://localhost:8080/events/query?param=value \
-H 'Content-Type: text/plain' \
-H "Authorization: Bearer <TOKEN>"
```
```bash
curl -X GET http://localhost:8080/events/query?customer_id=cus_118 \
-H 'Content-Type: text/plain' \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJ0b2tlbiI6ImV5SmhiR2NpT2lKSVV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9qRTJPVGcyT0RjeU9EbDkueFVnTGxOU0NNT0pzeEZKd1dEdmw4UmlNaE1hWTczdEZua2V6cDdiR2FjdyJ9.Qk1z90_ZMvf-l2IoTJ8rZIbh_uTLEbC3Fgj0aEmNHHU" 
```
```bash
# Example Response
[{"event_id":"652c96550ea634877babe4d9","event_type":"item_review","timestamp":"2023-10-16T01:48:05.073Z","service_name":"review_service","customer_id":"cus_118","event_details":{"action":"submitted_review","comment":"Very satisfactory!","item_id":"item892","rating":4}},
{"event_id":"652c978b0ea634877babe4da","event_type":"item_review","timestamp":"2023-10-16T01:53:15.164Z","service_name":"review_service","customer_id":"cus_118","event_details":{"action":"submitted_review","comment":"Bad","item_id":"item892","rating":0}}]
```
**Multiple query parameters are supported.**

## Schema 
```go
type Event struct {
	EventID      primitive.ObjectID     `json:"event_id" bson:"_id,omitempty"`
	EventType    string                 `json:"event_type" bson:"event_type"`
	Timestamp    time.Time              `json:"timestamp" bson:"timestamp"`
	ServiceName  string                 `json:"service_name" bson:"service_name"`
	CustomerID   string                 `json:"customer_id" bson:"customer_id"`
	EventDetails map[string]interface{} `json:"event_details" bson:"event_details"`
}
```
The **Event schema** is designed to capture audit log events with the following fields:

- **EventID:** A unique identifier for each event. Using MongoDB's ObjectID ensures uniqueness across any distributed scenario. **Filled automatically**
- **EventType:** Helps categorize the event. Useful in filtering and querying events of a certain type. e.g. **profile_updated**, **account_created**, **customer_billed**
- **Timestamp:** Automatically set to the current time when an event is created. Essential for audit logs to track when events occurred. **Filled automatically**
- **ServiceName:** Identifies which service in your infrastructure caused the event. e.g. **account_service**, **billing_service**, **resource_service**
- **CustomerID:** Helps tie an event to a particular customer/user. Useful for understanding user behavior and auditing specific user actions.
- **EventDetails:** A flexible map for any additional information about the event. This ensures the schema is extensible for future needs.

This schema is versatile and can be adapted to accommodate various event types and details.

- **Schema Flexibility:** MongoDB is schema-less. Given the nature of audit logs, over time, the data that needs to be logged might evolve. New types of events or additional fields might need to be logged. With MongoDB, this is easily possible without having to change the schema or run migrations.
- **Scalability:** MongoDB is designed to be scalable. It can handle huge amounts of data and offers horizontal scalability via sharding. This is particularly useful for logging systems which tend to generate large amounts of data over time.
- **High Write Throughput:** Audit logging systems typically have high write operations compared to read operations. MongoDB performs very well in scenarios with high write loads, thanks to its asynchronous write operations.
- **Indexing:** To make read operations faster, especially in scenarios where specific logs need to be retrieved based on certain criteria, MongoDB offers indexing. This ensures that logs can be queried efficiently even as the volume grows.



