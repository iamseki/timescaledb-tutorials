# Transactions on the Bitcoin blockchain

This project is an API that expose transactions on the Bitcoin blockchain data using timescaledb with go lang.

Covered concepts:

- Dependency injection using [fx](https://github.com/uber-go/fx)
- Tracing using [otel](https://opentelemetry.io/docs/instrumentation/go/)
- Querying analytical data using [timescaleDB](https://docs.timescale.com/api/latest/)
- Integrates CDC on a timescaleDB instance

## Initialize

- `docker compose up -d` && `./populate-db.sh`
- `go run main.go`

## Routes

| Endpoint           | Method        | Description                                                       | Response          |
| -------------------| --------------| ------------------------------------------------------------------| ------------------|
| `/v1/healthcheck`  | `GET`         | Basic healthcheck                                                 | `200` `500`       |
| `/v1/blocks/recent`| `GET`         | List recent blocks given a required limit (querystring)           | `200` `400` `500` |