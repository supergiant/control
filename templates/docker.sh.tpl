#!/bin/sh
sudo apt-get update \
  && sudo apt-get install -qy docker.io
sudo docker ps
