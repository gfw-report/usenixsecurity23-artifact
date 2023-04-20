#!/bin/bash

set -x
set -e

cd "$(dirname "$0")" || exit

servers=(
    eric6
    eric4
    eric5
    random5
    random6
    random7
    random8
    random9
    random10
)

for vm in "${servers[@]}"; do
    rsync -auvzP ../../affected-norand "$vm:~" &
done

wait

# stop previous servers
for vm in "${servers[@]}"; do
    ssh "$vm" 'tmux kill-server; tmux ls'  &
done

wait

# remove previous split files, as it may be from a different experiement.
for vm in "${servers[@]}"; do
    ssh "$vm" 'rm split*.ips*'  &
done

wait

rsync -auvP split_ambiguous-prefixes20.ips.aa eric6:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ab eric4:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ac eric5:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ad random5:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ae random6:~ &
rsync -auvP split_ambiguous-prefixes20.ips.af random7:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ag random8:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ah random9:~ &
rsync -auvP split_ambiguous-prefixes20.ips.ai random10:~ &

wait

ssh eric6 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-aa.pcap'" &
ssh eric4 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ab.pcap'" &
ssh eric5 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ac.pcap'" &
ssh random5 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ad.pcap'" &
ssh random6 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ae.pcap'" &
ssh random7 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-af.pcap'" &
ssh random8 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ag.pcap'" &
ssh random9 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ah.pcap'" &
ssh random10 "tmux new -d 'cd && sudo tcpdump -n not port 22 and tcp -w zero-ai.pcap'" &

wait

# # sleep 5

ssh eric6 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-aa.out'" &
ssh eric4 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ab.out'" &
ssh eric5 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ac.out'" &
ssh random5 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ad.out'" &
ssh random6 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ae.out'" &
ssh random7 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-af.out'" &
ssh random8 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ag.out'" &
ssh random9 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ah.out'" &
ssh random10 "tmux new -d 'cd && cat split* | ./affected-norand -try 3 -repeat 25 -worker 400 -p 80 -sleep 1s -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee zero-ai.out'" &

wait
