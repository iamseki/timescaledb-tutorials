# Analyze NYC taxi cab data

This project is a simple API to expose new york city taxi cab data using timescaledb with go programming language.

## Setup Environment :scroll:

- `docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg14-latest`

### Create database and tables
- `docker exec -ti timescaledb psql -U postgres -h localhost`
```SQL
CREATE DATABASE nyc_taxi_cab;
CREATE EXTENSION IF NOT EXISTS timescaledb;
```

- `docker cp tables.sql timescaledb:/home/postgres/tables.sql`
- `docker exec timescaledb psql -U postgres -h localhost -d nyc_taxi_cab -f tables.sql`
- test connection: `docker exec -ti timescaledb psql -U postgres -h localhost -d nyc_taxi_cab`

### Populate database

- Run the script `./populate-db.sh`
- This simple test must have an output: `docker exec timescaledb psql -U postgres -h localhost -d nyc_taxi_cab -c "SELECT * FROM rides LIMIT 5;"`

### Run API

- Just execute `go run main.go`

## Routes ::

- GET `http://localhost:5000/v1/rides/day/since?date=2023-01-01`
- GET `http://localhost:5000/v1/rides/fare/day/since?date=2023-01-01`