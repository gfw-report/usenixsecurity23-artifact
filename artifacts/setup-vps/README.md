# SSH configuration

1. Generate a pair of SSH key: `ssh-keygen -t ed25519 -f ~/.ssh/usenix-ae-reviewer -C ""` on your local machine.

2. Purchase two VPSes with the SSH public key installed: one inside China and the other outside of China.

3. Append the following text to `~/.ssh/config`:

```txt
Host usenix-ae-*
     Port 22
     User root
     IdentityFile ~/.ssh/usenix-ae-reviewer

Host usenix-ae-server-us
     HostName REDACTED_US_SERVER_IP

Host usenix-ae-client-china
     HostName REDACTED_CN_SERVER_IP
```

4. Now you should be able to log in to the VPSes with:

```sh
ssh usenix-ae-server-us
ssh usenix-ae-client-china
```

5. run this script to compile all necessary tools, and rsync the binaries to the server:

```sh
./setup-client/to_alibaba_server.sh
```

6. Run this script to compile a sink server, rsync the binary to the server, and make the server listen on all ports from 0-65535 for the specific client IP address:

```sh
./setup-server/to_digitalocean_server.sh
```


## Basic test

To quickly test if you have properly set up a pair of VPSes that can trigger the blocking:

Send some random probes from `usenix-ae-client-china` to port `2` of `usenix-ae-server-us` by repetitively executing the following command:

```sh
head -c200 /dev/urandom | nc -vn REDACTED_US_SERVER_IP 2
```

After executing the command a few times (1 time to 15 times), you will notice the `nc` cannot connect to `REDACTED_US_SERVER_IP:2` anymore. Congratulations! The blocking is triggered. You should still be able to connect to other ports of the same server, for example, `REDACTED_US_SERVER_IP:3`.

Alternatively, you can use the triggering tools:

```sh
echo REDACTED_US_SERVER_IP | ./utils/affected-norand -p 2 -log /dev/null
```

This tool will take a list of IPs on stdin, and perform (default 25) repeated connections to
the specified port, sending the
same (configurable) random payload in each connection. If the tool is unable to connect for
(default 5) consecutive connections in a row, the tool labels the IP as `affected` by
blocking (`true` in the `affected` column):

```txt
endTime,addr,countSuccess,totalTimeout,consecutiveTimeout,code,affected
1678218390,REDACTED_US_SERVER_IP:2,2,5,5,timeout,true
```
