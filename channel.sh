#!/bin/bash

channel_file=/home/raoni/tstando2

message=

for word in "$@"
do
  message="$message*$word*"
done

echo $message > $channel_file