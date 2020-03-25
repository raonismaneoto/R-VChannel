#!/bin/bash

sudo su

apt install ufw

port_to_open=$1

ufw allow $port_to_open/tcp

file_to_create=$2

touch $file_to_create
