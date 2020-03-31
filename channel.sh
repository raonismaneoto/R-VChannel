#!/bin/bash

channel_file=/home/raoni/tstando3

message=

for word in "$@"
do
  message="$message*$word*"
done

echo $message > $channel_file