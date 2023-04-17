#!/usr/bin/bash

export PGPASSWORD=password

# start postgres container
docker-compose up -d db

echo
echo waiting for the DB to start
echo
# sleep 10

until psql -U "postgres" -h 0.0.0.0 -p 5433 -c '\q'; do
	echo $PGPASSWORD_FILE
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

# get postgres container id
container_id=$(docker ps | grep ".*_db_1$" | awk '{print $1}')

cd ./portCaptureServer/db/

./setupDB.sh

# stop postgres container
docker stop $container_id
