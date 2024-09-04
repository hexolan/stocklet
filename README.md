# Stocklet

An event-driven microservices-based distributed e-commerce example application written in Golang. *(mouthful)*

## 📘 About

This project was originally built as an experiment with event-driven architecture. But I hope it can future serve as a beneficial demonstration of utilising the architecture and exemplify the implementation of some other miscellaneous microservice patterns.

Any ideas, suggestions or direct contributions to better conform with general and evolving industry practices are welcome and will be greatly appreciated, as I'd like for this project to evolve to the stage of being somewhat a reflection of a production-ready enterprise application.

⚠️ The application should be considered in the experimental prototype stage. Breaking changes can be expected between any future commits to this repo, in order to ease the development process and allow for clean refactoring of the project.

## 📝 Features

* Monorepository layout
* Microservice architecture
* Event-driven architecture
* Interfacing with services using gRPC
* User-facing RESTful HTTP APIs with gRPC-Gateway
* Distributed tracing with OpenTelemetry
* Transactional outbox pattern with Debezium
* API gateway pattern using Envoy
* Distributed transactions utilising the saga pattern

## ⚠️ Notice

As this project is licensed under the GNU Affero General Public License v3, [copying, templating or referencing code from this project](https://en.wikipedia.org/wiki/Clean-room_design) may violate international copyright law unless your project is using a compatible open-source license. Please ensure any implementation in your own projects is original and complies with applicable licenses and laws.

In the nature of open-source software, please consider contributing and giving back to the project to help make it better for the greater community, especially if you see it as a useful learning resource.

## 🗃️ Architecture

### 🔎 Overview

![Architecture Overview](/docs/imgs/overview.svg)

### 🧰 Technical Stack

#### Libraries, Frameworks and Tools

* API Tooling
  * [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc)
  * [github.com/grpc-ecosystem/grpc-gateway/v2](https://pkg.go.dev/github.com/grpc-ecosystem/grpc-gateway/v2)

* Client Libraries
  * [go.opentelemetry.io/otel](https://pkg.go.dev/go.opentelemetry.io/otel)
  * [github.com/twmb/franz-go](https://pkg.go.dev/github.com/twmb/franz-go)
  * [github.com/jackc/pgx/v5](https://pkg.go.dev/github.com/jackc/pgx/v5)

* Protobuf Libraries
  * [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf)
  * [github.com/bufbuild/protovalidate-go](https://pkg.go.dev/github.com/bufbuild/protovalidate-go)

* Tools
  * [plantuml.com](https://plantuml.com/)
  * [github.com/bufbuild/buf/cmd/buf](https://buf.build/docs/installation)
  * [github.com/golang-migrate/migrate/v4](https://pkg.go.dev/github.com/golang-migrate/migrate/v4#section-readme)

* Miscellaneous
  * [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
  * [github.com/rs/zerolog](https://pkg.go.dev/github.com/rs/zerolog)
  * [github.com/lestrrat-go/jwx/v2](https://pkg.go.dev/github.com/lestrrat-go/jwx/v2)
  * [github.com/doug-martin/goqu/v9](https://pkg.go.dev/github.com/doug-martin/goqu/v9)

#### Infrastructure

* Message Brokers
  * [Kafka](https://hub.docker.com/r/bitnami/kafka)
* Databases
  * [PostgreSQL](https://hub.docker.com/_/postgres)
* Miscellaneous
  * [OpenTelemetry](https://opentelemetry.io/)
  * [Envoy](https://www.envoyproxy.io/)
  * [Debezium Connect](https://hub.docker.com/r/debezium/connect)
* Provisioning and Deployment
  * [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### 🧩 Services

| Name | gRPC (w/ Gateway) | Produces Events | Consumes Events |
| :-: | :-: | :-: | :-: |
| [auth](/internal/svc/auth/) | ✔️ | ❌ | ✔️ |
| [order](/internal/svc/order/) | ✔️ | ✔️ | ✔️ |
| [payment](/internal/svc/payment/) | ✔️ | ✔️ | ✔️ |
| [product](/internal/svc/product/) | ✔️ | ✔️ | ✔️ |
| [shipping](/internal/svc/shipping/) | ✔️ | ✔️ | ✔️ |
| [user](/internal/svc/user/) | ✔️ | ✔️ | ❌ |
| [warehouse](/internal/svc/warehouse/) | ✔️ | ✔️ | ✔️ |

Each service is prepared by a [``service-init``](/cmd/service-init/) container; a deployment responsible for performing any database migrations and configuring Debezium outbox connectors for that service.

### 📇 Events

The events are schemed and serialised using [protocol buffers](https://protobuf.dev/). They are dispatched using the [transactional outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html), with [Debezium](https://debezium.io/) used as a relay to read and publish events from database outbox tables to the message broker.

Further documentation on the events can be found at [``/docs/EVENTS.md``](/docs/EVENTS.md)

## 💻 Deployment

### Using Docker

The application can be deployed using [Docker Compose](https://docs.docker.com/compose/) (with the compose files located in [``/deploy/docker/``](/deploy/docker/)). Ensure the correct configuration is in place by copying and removing ``.example`` from the end of the example environment files located in [``/deploy/configs/``](/deploy/configs/).

Deploy using the following command: ``docker compose -f deploy/docker/compose.yaml -f deploy/docker/compose.override.yaml up --build``

## 🧪 Contributing

If you like this project then please leave a ⭐ to show your support. All forms of feedback and contributions are welcome and greatly appreciated!

Have any [ideas for improvements?](/docs/ROADMAP.md) Please don't hesistate to [open an issue](https://github.com/hexolan/stocklet/issues/new) to discuss, or a [pull request](https://github.com/hexolan/stocklet/compare) with [enhancements](https://github.com/hexolan/stocklet/fork).

## 📓 License

This project is licensed under the [GNU Affero General Public License v3](/LICENSE).
