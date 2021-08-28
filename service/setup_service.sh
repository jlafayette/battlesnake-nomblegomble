#!/usr/bin/env bash
set -eux

echo "Running setup_services.sh"

echo "create /root/bin directory"
mkdir -p /root/bin

echo "copy service file"
mv ./battlesnake.service /etc/systemd/system/battlesnake.service

echo "start service"
systemctl start battlesnake

echo "Done!"