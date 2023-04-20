#!/bin/bash

set -x
set -e

cd "$(dirname "$0")" || exit

names=(
    aa
    ab
    ac
    ad
    ae
    af
    ag
    ah
    ai
)

for name in "${names[@]}"; do
    tshark -r "${name}.pcap" -Y "(tcp.flags == 0x012) && (tcp.window_size_value == 0)" -Tfields -E separator=\; -e ip.src | sort | uniq > "win0.${name}.ips" &
done

wait

cat win0.a*.ips | sort | uniq > aggregated_win0.ips
