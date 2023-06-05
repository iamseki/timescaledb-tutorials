# Analyze NYC taxi cab data

## Setup Environment :scroll:

- `docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg14-latest`

### Create database and tables
- `docker exec -ti timescaledb psql -U postgres -h localhost`
```SQL
CREATE DATABASE nyc_taxi_cab;
CREATE EXTENSION IF NOT EXISTS timescaledb;
```

tables: 
- `docker cp tables.sql timescaledb:/home/postgres/tables.sql`
- `docker exec timescaledb psql -U postgres -h localhost -d nyc_taxi_cab -f tables.sql`
- test connection: `docker exec -ti timescaledb psql -U postgres -h localhost -d nyc_taxi_cab`

### Populate database

- Run the script `./populate-db.sh`
- This simple test must have an output: `docker exec timescaledb psql -U postgres -h localhost -d nyc_taxi_cab -c "SELECT * FROM rides LIMIT 5;"`