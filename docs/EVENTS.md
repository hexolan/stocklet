# Stocklet Docs: Events

## Table of Contents

* [Repository Overview](/README.md)
* [Documentation: Overview](/docs/README.md)
* [Documentation: Events](/docs/EVENTS.md)
* [Documentation: Feature Roadmap](/docs/ROADMAP.md)

## Overview

The events are schemed and serialised using [protocol buffers](https://protobuf.dev/). The events schemas can be found in [``/schema/protobufs/stocklet/events/``](/schema/protobufs/stocklet/events/)

They are dispatched using the [transactional outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html). Debezium is used as a relay to publish events from database outbox tables to the message broker (Kafka). The Debezium connectors are configured by the ``service-init`` containers, which are also responsible for performing database migrations for their respective services.

## Services

### Auth Service

**Produces:**

* n/a

**Consumes:**

* UserDeletedEvent

### Order Service

**Produces:**

* OrderCreatedEvent
* OrderPendingEvent
* OrderRejectedEvent
* OrderApprovedEvent

**Consumes:**

* ProductPriceQuoteEvent
* StockReservationEvent
* ShipmentAllocationEvent
* PaymentProcessedEvent

### Payment Service

**Produces:**

* BalanceCreatedEvent
* BalanceCreditedEvent
* BalanceDebitedEvent
* BalanceClosedEvent
* TransactionLoggedEvent
* TransactionReversedEvent *(currently unused)*
* PaymentProcessedEvent

**Consumes:**

* UserCreatedEvent
* UserDeletedEvent
* ShipmentAllocationEvent

### Product Service

**Produces:**

* ProductCreatedEvent
* ProductPriceUpdatedEvent
* ProductDeletedEvent
* ProductPriceQuoteEvent

**Consumes:**

* OrderCreatedEvent

### Shipping Service

**Produces:**

* ShipmentAllocationEvent
* ShipmentDispatchedEvent *(currently unused)*

**Consumes:**

* StockReservationEvent
* PaymentProcessedEvent

### User Service

**Produces:**

* UserCreatedEvent
* UserEmailUpdatedEvent
* UserDeletedEvent

**Consumes:**

* n/a

### Warehouse Service

**Produces:**

* StockCreatedEvent
* StockAddedEvent *(currently unused)*
* StockRemovedEvent
* StockReservationEvent

**Consumes:**

* OrderPendingEvent
* ShipmentAllocationEvent
* PaymentProcessedEvent

## Miscellaneous

### Place Order Saga

The place order [saga](https://microservices.io/patterns/data/saga.html) is initiated when a new order is created.

![Place Order Saga](/docs/imgs/placeordersaga.svg)
