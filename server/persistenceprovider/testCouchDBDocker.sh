#!/usr/bin/env bash

echo "### - Testing if the container is already running..."

docker ps | grep my_couchdb 1>/dev/null

if [ $? -eq 0 ] ; then
  echo -e "\n################ Stoping running container ###################\n"
  docker stop my_couchdb
fi

echo "### - Testing if the container already exists..."
docker ps -a | grep my_couchdb 1>/dev/null

if [ $? -eq 0 ] ; then
  echo -e "\n################ Removing existing container ##################\n"
  docker rm my_couchdb
fi

set -e

docker run --name my_couchdb -p 5984:5984 -d couchdb 1>/dev/null
go test -v -coverprofile=profile.out -covermode=atomic -tags integration
if [ -f profile.out ]; then
    cat profile.out >> ../coverage.txt
    rm profile.out
fi
docker stop my_couchdb 1>/dev/null
docker rm my_couchdb 1>/dev/null
