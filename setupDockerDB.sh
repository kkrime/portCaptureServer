#!/usr/bin/bash

# start postgres container
docker-compose up -d db

echo
echo waiting for the DB to start
echo
sleep 10

# get postgres container id
container_id=$(docker ps | grep ".*_db_1$" | awk '{print $1}')

cd ./portCaptureServer/db/

./setupDB.sh

# stop postgres container
docker stop $container_id
