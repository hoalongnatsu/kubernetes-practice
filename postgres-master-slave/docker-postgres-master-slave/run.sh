#!/bin/bash
opt=$1

if [[ $opt = "up" ]]; then
  docker-compose up -d
elif [[ $opt = "down" ]]; then
  docker-compose down
fi
