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

## Querying :scroll:

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

## Compression :small_blue_diamond:

To enable compression uses `timescaledb.compress`, the orderby and segmentby will be default to time if not specified, more details on [compression docs](https://docs.timescale.com/use-timescale/latest/compression/).

```SQL
ALTER TABLE stocks_real_time SET (
  timescaledb.compress, 
  timescaledb.compress_orderby = 'time DESC', 
  timescaledb.compress_segmentby = 'symbol'
);

--- to check compression_setings:
SELECT * FROM timescaledb_information.compression_settings;
```

### Automatic Compression

We can schedule a policy to [automatically compress](https://docs.timescale.com/api/latest/compression/add_compression_policy/) the data. For example, if you want to compress hypertable data that is older than two weeks, run:

```SQL
SELECT add_compression_policy('stocks_real_time', INTERVAL '2 weeks');
--- TO SEE POLICY DETAILS:
SELECT * FROM timescaledb_information.jobs;
--- TO see Job Statistics:
SELECT * FROM timescaledb_information.job_stats;
```

Compressed rows can't be updated or deleted!!1! So it's bets to compress aged or data is less likely to require updating..

### Manual Compression

When [compressing manually chunks](https://docs.timescale.com/api/latest/compression/compress_chunk/) is needed.

An example to compress chunks that consist of data older than 2 weeks:

```SQL
SELECT compress_chunk(i, if_not_compressed=>true) -- if statement to not shows an error when it tries to compress an already compressed chunk!
  FROM show_chunks('stocks_real_time', older_than => INTERVAL '2 weeks') i;
```

### Verifying Compression

Check queries:

```SQL
SELECT 
  pg_size_pretty(before_compression_total_bytes) as "before compression", 
  pg_size_pretty(after_compression_total_bytes) as "after compression"
FROM hypertable_compression_stats('stocks_real_time');
```

## Data Retention :clock1:

> An intrisinc part of working with time-series data is that the relevance of data can diminish over time.

This is a feature that I'm not interested but know it's possible to use [data retention policy](https://docs.timescale.com/getting-started/latest/data-retention/#create-a-data-retention-policy).