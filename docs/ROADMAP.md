# Stocklet Docs: Feature Roadmap

## Table of Contents

* [Repository Overview](/README.md)
* [Documentation: Overview](/docs/README.md)
* [Documentation: Events](/docs/EVENTS.md)
* [Documentation: Feature Roadmap](/docs/ROADMAP.md)

## Prologue

This document should be considered a brainstorm of ideas.

There is no guarantee that I will implement any listed features or functionality below. After all, I initially made this application as an experiment with EDA, and there are areas that could use improvement and expansion (the application is a prototype in current format). Some of the current implemented functionality is quite bare-bones, so if I come to revisit this project at a later date this document is where I'd first look.

However, contributions are welcome; if you feel like implementing something (already below or not), or otherwise spot other areas that could use improvement, then please feel free to open an issue to discuss or a pull request with your implementation.

## Feature Ideas

* Front-end user interface
  * Allow interfaces with the application through alternatives means
* Notification service
  * Send notifications to users (e.g. through a mock email or a unread messages mechanism) upon reciept of events related to order status changes (i.e. OrderApprovedEvent)
* Product recommendation service
  * Provide a list of recommended products catered to specific customers

## Miscellaneous Ideas

* Integration tests
* Ensured idempotency in event consumers (service-side)
* Clear-up of event processes
* Kubernetes deployment (prepare manifest files)
* Interchangable infrastructure
  * Support for NATS as a message broker
  * Support for MongoDB as a database
