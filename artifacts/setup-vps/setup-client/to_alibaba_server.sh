#!/bin/bash

set -x
set -e

cd "$(dirname "$0")" || exit

servers=(
    usenix-ae-client-china
)

# compile utils
cd ../../utils && make && cd -

for vm in "${servers[@]}"; do
    rsync -auvzP alibabacloud.sh sysctl.conf ../../utils ../../test-* "$vm:~" &
done

wait

# make sure all servers are reachable
for vm in "${servers[@]}"; do
    ssh "$vm" 'tmux kill-server; tmux ls'  &
done

wait

sleep 5

for vm in "${servers[@]}"; do
    ssh "$vm" 'bash ~/alibabacloud.sh'  &
done

wait
