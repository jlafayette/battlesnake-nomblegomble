#!/usr/bin/env bash
set -eux

echo "Running update.sh script!"

echo "stop service"
systemctl stop battlesnake

echo "copy executable file"
chmod +x ./battlesnake-go
mv ./battlesnake-go /root/bin/battlesnake-go

echo "start service"
systemctl start battlesnake

echo "Done!"
