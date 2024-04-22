#!/bin/bash

kind delete cluster &

docker image rm http_red:v0.0.1 controller:v0.0.1 example_workload:v0.0.1 load_generator:v0.0.1 &

docker image prune -f &

for job in `jobs -p`
do
    wait $job
done
