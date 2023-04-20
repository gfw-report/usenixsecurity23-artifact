#!/usr/bin/env python3

random_payload = "dadd034913c52da75fd9f05dc76803917134808efed97ef8884f2151b712f60fed634f609f132033a15b77ed3ccaa2d20f5b"

for j in range(5, 7):
    for i in range(256):
        print(hex(i)[2:].zfill(2) * j + random_payload)
