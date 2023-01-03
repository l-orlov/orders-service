# orders-service

Service for working with orders data.  
You can interact with it via the http api and also through a channel in the NATS Streaming.  
PostgreSQL database and local cache are used to store data.

### Local running

* Set up dependencies in docker:  
  `make docker-up-local`
* Set up migrations in db:  
  `make db-reset-local`
* Start service:  
  `make run`
* Send message with test order to nats-streaming:  
  `make send-test-message`

### HTTP API

* `GET http://localhost:8080/orders` - get all orders
* `GET http://localhost:8080/orders/:id` - get order by id
* `POST http://localhost:8080/orders` - create order
