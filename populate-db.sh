FILE="$(pwd)"/real_time_stock_data.zip

if test -f "$FILE"; then
  echo "dataset file: already $FILE exists."
else 
  echo "dataset file $FILE does not exists, downloading..."
  curl https://assets.timescale.com/docs/downloads/get-started/real_time_stock_data.zip -o real_time_stock_data.zip
fi

UNZIP_FILE="$(pwd)"/tutorial_sample_tick.csv
if test -f "$UNZIP_FILE"; then
  echo "already unziped real_time_stock_data.zip"
else 
  echo "unziping real_time_stock_data.zip..."
  unzip real_time_stock_data.zip
fi

docker cp tutorial_sample_tick.csv timescaledb:/home/postgres/tutorial_sample_tick.csv
docker cp tutorial_sample_company.csv timescaledb:/home/postgres/tutorial_sample_company.csv

docker exec timescaledb psql -U postgres -h localhost -d stock_exchange -c "\COPY stocks_real_time from './tutorial_sample_tick.csv' DELIMITER ',' CSV HEADER;"
docker exec timescaledb psql -U postgres -h localhost -d stock_exchange -c "\COPY company from './tutorial_sample_company.csv' DELIMITER ',' CSV HEADER;"


