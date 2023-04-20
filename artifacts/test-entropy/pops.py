import random

N = 50*8 # length in bits

def bits_to_bytes(s):
    return int(s, 2).to_bytes((len(s)+7) // 8, byteorder='big')

def shuffle(s):
    l = list(s)
    random.shuffle(l)
    return ''.join(l)

def popcount(b):
    t = 0
    for c in b:
        t += sum([int(x) for x in bin(c)[2:]])
    return t

# i is the number of 1s in the string
for i in range(0, N):
    s = '1'*i + '0'*(N-i)
    b = bits_to_bytes(shuffle(s))
    #print(b.hex(), popcount(b))
    print(b.hex())
