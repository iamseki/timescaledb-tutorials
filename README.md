# Timescale DB - Get Started

## Setup Environment :hammer:

<details>
  <summary>SETUP</summary>

## Running timescaledb container

- `docker run -d --name timescaledb -p 5432:5432 -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg14-latest`
  - Check for an earlier version on [timescale-dockerhub](https://hub.docker.com/r/timescale/timescaledb)

- Connects to psql inside the container: `docker exec -it timescaledb psql -U postgres`

  - `CREATE DATABASE stock_exchange;`
  - `CREATE EXTENSION IF NOT EXISTS timescaledb;`
  - Connects to specific db: `docker exec -it timescaledb psql -U postgres -h localhost -d stock_exchange`
  - To list installed extensions: `\dx`

## Creating the Hypertable

> PGsql tabels that are automatically partitioning by time. more on [hypertables official docs](https://docs.timescale.com/use-timescale/latest/hypertables/)

- Creates a regular postgresql table:
  ```SQL
    CREATE TABLE stocks_real_time (
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    price DOUBLE PRECISION NULL,
    day_volume INT
  );
  ```
- Then the hyper table with:
  ```SQL
    SELECT create_hypertable('stocks_real_time', 'time'); --- (table, partition_key)
  ```

- Also creates an index:
  ```SQL
    CREATE INDEX idx_symbol_time ON stocks_real_time (symbol, time DESC);
  ```

## Enhance time-series data with an regular pgsql table and the hypertable

```SQL
  CREATE TABLE company (
    symbol TEXT NOT NULL,
    name TEXT NOT NULL
  );
```

## Adds the dataset

- This dataset its providen by Twelve Data and has ~8 milion rows.
- Just run the `./populate-db.sh` script. Remember the container named`timescaledb` should exists been up and running !1!.

</details>

## Querying

> Remember to be connected into container psql or some other pgsql client !!! `docker exec -it timescaledb psql -U postgres -h localhost -d stock_exchange`

The well known SQL works:


```SQL
  -- the most 10 recents samples:
  select * from stocks_real_time order by time desc limit 10;
  -- built in postgresql functions to know the periods of the data from this dataset
  select time::date as iso_date from stocks_real_time s group by s.time::date; 
```