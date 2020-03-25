#!/bin/bash

sudo apt install ufw

port_to_open=$1

sudo ufw allow $port_to_open/tcp

file_to_create=$2
touch $file_to_create

channel_file_pattern="channel_file="

sed -i "s,channel_file=.*,channel_file=$file_to_create," "$(pwd)/channel.sh"

exit $?
