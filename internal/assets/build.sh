#!/bin/bash

docker build -t docker-packr .
docker run -v $(pwd):/home/migrations docker-packr
