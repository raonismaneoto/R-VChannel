#!/bin/bash

sudo cp $0 /etc/init.d/

project_path="$HOME/go/src/github.com/raonismaneoto/R-VChannel/"

go build $project_path/RVC.go
sudo go run $project_path/RVC.go &

pid=$!

while :
do
  if ps -p $pid > /dev/null
  then
    sleep 300
    continue
  else
    $(pwd)/install.sh
  fi
done