#!/bin/bash

cd "$(dirname "$0")" || exit

# ssh security
echo -e "PermitRootLogin yes\nAllowUsers root\nPasswordAuthentication no" | sudo tee /etc/ssh/sshd_config.d/noroot.conf

# unattended security upgrade
sudo apt-get update && sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y
sudo apt-get install -y unattended-upgrades
sudo dpkg-reconfigure -f noninteractive -p low unattended-upgrades

# setup sysctl.conf
sudo cp sysctl.conf /etc/sysctl.d/10-sink.conf
# permanently increase max number of open files; it requires reboot
echo -e "* soft nofile 1048576\n* hard nofile 1048576\n* soft nproc 10240\n* hard nproc 10240" | sudo tee /etc/security/limits.d/99-sink.conf
