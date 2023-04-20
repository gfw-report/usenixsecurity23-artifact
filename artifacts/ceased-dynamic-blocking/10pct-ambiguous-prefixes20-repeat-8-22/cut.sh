#!/bin/bash

cd "$(dirname "$0")" || exit

(awk '{if ($5>10 && $3<0.8 && $3>0.2) print $1}' prefixes24.out | xargs nmap -sL -n) | grep 'Nmap scan report for' | cut -f 5 -d ' '  | shuf > ambiguous-prefixes20.ips
split -l 47000 ambiguous-prefixes20.ips split_ambiguous-prefixes20.ips.
