ALTER SYSTEM SET wal_level = logical;

CREATE TABLE IF NOT EXISTS transactions (
   time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   block_id INT,
   hash TEXT,
   size INT,
   weight INT,
   is_coinbase BOOLEAN,
   output_total BIGINT,
   output_total_usd DOUBLE PRECISION,
   fee BIGINT,
   fee_usd DOUBLE PRECISION,
   details JSONB
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS outbox_events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  aggregatetype VARCHAR(255) NOT NULL CHECK (
    aggregatetype IN (
    'INSERTED',
    'UPDATED',
    'DELETED'
  )),
  aggregateid VARCHAR(255) NOT NULL,
  payload JSONB NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW()  
);

SELECT create_hypertable('transactions', 'time');

CREATE INDEX IF NOT EXISTS hash_idx ON public.transactions USING HASH (hash);

CREATE INDEX IF NOT EXISTS block_idx ON public.transactions (block_id);

CREATE UNIQUE INDEX IF NOT EXISTS time_hash_idx ON public.transactions (time, hash);

---- Transactions info by hour with continuous aggregates

CREATE MATERIALIZED VIEW IF NOT EXISTS one_hour_transactions
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', time) AS bucket,
   count(*) AS tx_count,
   sum(fee) AS total_fee_sat,
   sum(fee_usd) AS total_fee_usd,
   stats_agg(fee) AS stats_fee_sat,
   avg(size) AS avg_tx_size,
   avg(weight) AS avg_tx_weight,
   count(
         CASE
            WHEN (fee > output_total) THEN hash
            ELSE NULL
         END) AS high_fee_count
  FROM transactions
  WHERE (is_coinbase IS NOT TRUE)
GROUP BY bucket;

-- Refresh policy
SELECT add_continuous_aggregate_policy('one_hour_transactions',
   start_offset => INTERVAL '3 hours',
   end_offset => INTERVAL '1 hour',
   schedule_interval => INTERVAL '1 hour');

---- Blocks info by hour with continuous aggregates
CREATE MATERIALIZED VIEW IF NOT EXISTS one_hour_blocks
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', time) AS bucket,
   block_id,
   count(*) AS tx_count,
   sum(fee) AS block_fee_sat,
   sum(fee_usd) AS block_fee_usd,
   stats_agg(fee) AS stats_tx_fee_sat,
   avg(size) AS avg_tx_size,
   avg(weight) AS avg_tx_weight,
   sum(size) AS block_size,
   sum(weight) AS block_weight,
   max(size) AS max_tx_size,
   max(weight) AS max_tx_weight,
   min(size) AS min_tx_size,
   min(weight) AS min_tx_weight
FROM transactions
WHERE is_coinbase IS NOT TRUE
GROUP BY bucket, block_id;

-- Refresh policy
SELECT add_continuous_aggregate_policy('one_hour_blocks',
   start_offset => INTERVAL '3 hours',
   end_offset => INTERVAL '1 hour',
   schedule_interval => INTERVAL '1 hour');


---- Transactions that miners received as rewards each hour
CREATE MATERIALIZED VIEW IF NOT EXISTS one_hour_coinbase
WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', time) AS bucket,
   count(*) AS tx_count,
   stats_agg(output_total, output_total_usd) AS stats_miner_revenue,
   min(output_total) AS min_miner_revenue,
   max(output_total) AS max_miner_revenue
FROM transactions
WHERE is_coinbase IS TRUE
GROUP BY bucket;

-- Refresh policy
SELECT add_continuous_aggregate_policy('one_hour_coinbase',
   start_offset => INTERVAL '3 hours',
   end_offset => INTERVAL '1 hour',
   schedule_interval => INTERVAL '1 hour');