docker cp tables.sql timescaledb:/home/postgres/tables.sql

docker exec timescaledb psql -U postgres -h localhost -d nyc_taxi_cab -f tables.sql

FILE="$(pwd)"/nyc_data.tar.gz

if test -f "$FILE"; then
  echo "dataset file: already $FILE exists."
else 
  echo "dataset file $FILE does not exists, downloading..."
  curl https://assets.timescale.com/docs/downloads/nyc_data.tar.gz -o nyc_data.tar.gz
fi

UNZIP_FILE="$(pwd)"/nyc_data_rides.csv
if test -f "$UNZIP_FILE"; then
  echo "already unziped nyc_data.tar.gz"
else 
  echo "unziping nyc_data.tar.gz..."
  tar -xvf nyc_data.tar.gz
fi

docker cp nyc_data_rides.csv timescaledb:/home/postgres/nyc_data_rides.csv

docker exec timescaledb psql -U postgres -h localhost -d nyc_taxi_cab -c "\COPY rides from './nyc_data_rides.csv' DELIMITER ',' CSV HEADER;"
