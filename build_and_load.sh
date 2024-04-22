#!/bin/bash

docker build . -t $1 ; kind load docker-image $1
