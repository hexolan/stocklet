# Stocklet Docs

## Table of Contents

* [Repository Overview](/README.md)
* [Documentation: Overview](/docs/README.md)
* [Documentation: Events](/docs/EVENTS.md)
* [Documentation: Feature Roadmap](/docs/ROADMAP.md)

## Formatting, styling, etc

The code has been formatted using [`gofmt`](https://pkg.go.dev/cmd/gofmt).

The protobuf schema files have been formatted using [`buf format`](https://buf.build/docs/reference/cli/buf/format).

The markdown files have been linted and formatted using [markdownlint](https://github.com/DavidAnson/markdownlint) (with the exception of MD013).

Commit messages should adhere to the [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification.

I used [PlantUML](https://plantuml.com/) as the tool to help make the diagrams. The PlantUML files are available alongside the resulting images: [`/docs/imgs/overview.plantuml`](/docs/imgs/overview.plantuml) and [`/docs/imgs/placeordersaga.plantuml`](/docs/imgs/placeordersaga.plantuml)
