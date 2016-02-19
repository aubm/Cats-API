#!/bin/bash

docker build -t kendo5731/cats-api .
docker-compose up -d --force-recreate
docker ps -a
go test ./...
newman -c cats_api.json.postman_collection -e cats_api.postman_environment -x
