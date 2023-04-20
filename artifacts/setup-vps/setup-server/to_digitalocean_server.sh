#!/bin/bash

set -x
set -e

cd "$(dirname "$0")" || exit

servers=(
    usenix-ae-server-us
)

# only once: make sure no tmux running
for vm in "${servers[@]}"; do
    ssh "$vm" 'tmux kill-server; tmux ls; sudo pkill sink'  &
done
wait

# only once: configure VPSes
for vm in "${servers[@]}"; do
    (
	rsync -auvP digitalocean.sh sysctl.conf "$vm:~"
	ssh "$vm" "tmux new -d 'bash -l ~/digitalocean.sh'"
    ) &
done
wait

# only once: load iptables rules
for vm in "${servers[@]}"; do
    rsync -auvP iptables "$vm:~" &
    ssh "$vm" "tmux new -d 'bash -l ~/iptables/iptables_only_targeted_sip.sh'" &
    # ssh "$vm" "tmux new -d 'sudo iptables-restore < ./iptables/iptables_rules'" &
done
wait

# make sure no tmux running
for vm in "${servers[@]}"; do
    ssh "$vm" 'tmux kill-server; tmux ls; sudo pkill sink' &
done
wait

# compile sink server
cd sink-server && env GOOS=linux GOARCH=amd64 go build && cd -

for vm in "${servers[@]}"; do
    (
	rsync -auvzP sink-server/sink "$vm:~"
	ssh "$vm" "tmux new -d './sink'" &
    ) &
done
wait

sleep 5
