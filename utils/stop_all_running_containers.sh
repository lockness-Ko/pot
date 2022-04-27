#! /bin/bash

docker ps -q | while read -r line; do docker rm -f $line; done
