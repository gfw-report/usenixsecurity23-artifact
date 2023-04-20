#!/bin/bash

cd "$(dirname "$0")" || exit

set -x

# backup
sudo iptables-save > "iptables-rules-backup-$(date '+%Y-%m-%d-%H:%M:%S').v4"

testIPs=(
    1.1.1.1
    8.8.8.8
)

for testIP in "${testIPs[@]}"; do
    # redirect all tcp packets from testIP to server's port 0-65535 to a single port 12345
    sudo iptables -t nat -A PREROUTING -i eth0 -p tcp -s "$testIP" --dport 0:65535 -j REDIRECT --to-port 12345
    # drop packets with RST flag set, regardless of other TCP flags
    iptables -I OUTPUT -p tcp -d "$testIP" --tcp-flags RST RST -j DROP
done

# not targeted testIP specific:
# drop any response from sport 53; note that we do not drop incoming
# packets to port 53 so that we can still inform whether packets are
# still being sent to the port 53.
iptables -I OUTPUT -p udp --sport 53 -j DROP
iptables -I OUTPUT -p icmp --icmp-type destination-unreachable -j DROP
iptables -I OUTPUT -p icmp --icmp-type port-unreachable -j DROP
