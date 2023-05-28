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

## Querying :scrol:

> Remember to be connected into container psql or some other pgsql client !!! `docker exec -it timescaledb psql -U postgres -h localhost -d stock_exchange`

### The well known SQL still works


```SQL
  -- the most 10 recents samples:
  SELECT * FROM stocks_real_time ORDER BY time DESC LIMIT 10;
  -- built in postgresql functions to know the periods of the data from this dataset
  SELECT time::date as iso_date FROM stocks_real_time s GROUP BY s.time::date;
  -- from the previous four days
  SELECT * FROM stocks_real_time srt WHERE time > now() - INTERVAL '4 days' LIMIT 10;

  -- avg trade price for Apple from the last four days
  SELECT
    avg(price)
  FROM stocks_real_time srt
  JOIN company c ON c.symbol = srt.symbol
  WHERE c.name = 'Apple' AND time > now() - INTERVAL '4 days';
```

### Timescale SQL functions

- There are a lot of [custom sql functions](https://docs.timescale.com/api/latest/hyperfunctions/) to perform time-series analysis.

- Using [first](https://docs.timescale.com/api/latest/hyperfunctions/first/) and [last](https://docs.timescale.com/api/latest/hyperfunctions/last/) to find the first and last trading price:

```SQL
  SELECT symbol, first(price, time), last(price, time)
  FROM stocks_real_time srt
  WHERE time > now () - INTERVAL '3 days'
  GROUP BY symbol -- since it is an aggregation function its needs a "GROUP BY" statement
  ORDER BY symbol;
```

- Using `time_bucket()` to bucket values based on an interval. More on [docs](https://docs.timescale.com/api/latest/hyperfunctions/time_bucket/). In this case to calculate the average daily price of each trading symbol over the last week:

```SQL
  SELECT time_bucket('1 day', time) AS bucket, symbol, avg(price)
  FROM stocks_real_time srt 
  WHERE time > now() - INTERVAL '1 week'
  GROUP BY bucket, symbol
  ORDER BY bucket, symbol
  LIMIT 10;
```

## Continuous Aggregate :pencil2:

Calculating aggregates(average price per day, maximum CPU last 5 minutes and so on) can be computationally intensive. Some reasons for that:

- Aggregating _large amounts of data_ often requires a lot of calculation time.
- Ingesting new data requires new aggregation calculations which can effect `ingest rate` and `aggregation speed`.

Timescale [continuous aggregates](https://docs.timescale.com/use-timescale/latest/continuous-aggregates/) solve both of theses problems. In summary the feature is a automatically refreshed materialized view.

### Use case

Lets use an aggregation query to generates [candlestick](https://en.wikipedia.org/wiki/Candlestick_chart) data that can be used to show a candlestick chart, calculating the `high`, `open`, `close` and `low` for prices given an interval:

```SQL
  SELECT
    time_bucket('1 day', time) AS day,
    symbol,
    max(price) AS high,
    first(price, time) AS open,
    last(price, time) AS close,
    min(price) AS low
  FROM stocks_real_time srt
  GROUP BY day, symbol
  ORDER BY day DESC, symbol;
```

Now we can creates a continuoues aggregate:

```SQL
  CREATE MATERIALIZED VIEW stock_candlestick_daily
  WITH (timescaledb.continuous) AS 
  SELECT
    time_bucket('1 day', time) AS day,
    symbol,
    max(price) AS high,
    first(price, time) AS open,
    last(price, time) AS close,
    min(price) AS low
  FROM stocks_real_time srt
  GROUP BY day, symbol;
```

After a while we can query the mv:

```SQL
  SELECT * FROM stock_candlestick_daily
  ORDER BY day DESC, symbol;
```

To inspect details about continuous aggregate:

```SQL
  SELECT * FROM timescaledb_information.continuous_aggregates;
```

### Policy

In Timescale 1.7 and later, real time aggregates are enabled by default, always including the most recent data (respect the next scheduled refresh)

We can set an automatic refresh [policy](https://docs.timescale.com/api/latest/continuous-aggregates/add_continuous_aggregate_policy/) to update the contunuous aggregate:

```SQL
  SELECT add_continuous_aggregate_policy('stock_candlestick_daily',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    scheduled_interval => INTERVAL '1 days'
    );

  --- to manually update given a period
  CALL refresh_continuous_aggregate(
    'stock_candlestick_daily',
    now() - INTERVAL '1 week',
    now()
  );
```

This policy runs once a day. When it runs, it materializes data from between 3 days ago and 1 hour ago.