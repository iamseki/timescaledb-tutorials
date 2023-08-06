ALTER SYSTEM SET wal_level = logical;

CREATE TABLE transactions (
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

SELECT create_hypertable('transactions', 'time');

CREATE INDEX hash_idx ON public.transactions USING HASH (hash);

CREATE INDEX block_idx ON public.transactions (block_id);

CREATE UNIQUE INDEX time_hash_idx ON public.transactions (time, hash);
