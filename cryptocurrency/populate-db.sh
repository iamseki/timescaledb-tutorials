docker cp database_up.sql cryptocurrency-timescaledb-1:/database_up.sql

docker exec cryptocurrency-timescaledb-1 psql -U postgres -h localhost -f '/database_up.sql'

FILE="$(pwd)"/bitcoin_sample.zip

if test -f "$FILE"; then
  echo "dataset file: already $FILE exists."
else 
  echo "dataset file $FILE does not exists, downloading..."
  curl https://assets.timescale.com/docs/downloads/bitcoin-blockchain/bitcoin_sample.zip -o bitcoin_sample.zip
fi

UNZIP_FILE="$(pwd)"/tutorial_bitcoin_sample.csv
if test -f "$UNZIP_FILE"; then
  echo "already unziped bitcoin_sample.zip"
else 
  echo "unziping bitcoin_sample.zip..."
  unzip bitcoin_sample.zip
fi

docker cp tutorial_bitcoin_sample.csv cryptocurrency-timescaledb-1:/tutorial_bitcoin_sample.csv

docker exec cryptocurrency-timescaledb-1 psql -U postgres -h localhost -c "\COPY transactions from './tutorial_bitcoin_sample.csv' DELIMITER ',' CSV HEADER;"