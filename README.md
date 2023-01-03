# orders-service

Service for working with orders data.  
You can interact with it via the http api and also through a channel in the NATS Streaming.  
PostgreSQL database and local cache are used to store data.

### Local running
Start service:  
`make run`  
Send message with test order to nats-streaming:  
`make send-test-message`  

### HTTP API

* GET /orders - get all orders
* GET /orders/:id - get order by id
* POST /orders - create order
