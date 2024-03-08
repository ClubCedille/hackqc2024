# Hackaton template

## Dev dependencies

- [Go 1.22](https://go.dev/doc/install)
- [Air](https://github.com/cosmtrek/air)

## Setup

### Nix

You can use nix to setup dependencies: `nix develop`

### Using Make (alternative)

Alternatively, you can manually initialize the project using the Makefile : `make init`

## How to run

In project root, run one of the following commands:

Go run:

`go run .`

Live reload:

`air`

### Using Makefile

```bash
make build
make start
```

#### with Docker

 Build the Docker image and run the container on port 8080

```bash
make docker-build
make docker-run
```
