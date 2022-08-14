# fwdctl

![fwdctl](./fwdctl.png)

**fwdctl** is a simple and intuitive CLI to manage forwards in your **Linux** server.

## How it works

It essentially provides commands to manage forwards, using **iptables** under the hood.

Let's do an example:

Suppose you have an **hypervisor** server that hosts some virtual machines inside itself. If you need to expose an internal service, managed by one of these VMs, you can use **fwdctl** from the hypervisor to add the forward to expose this service.

![example](./fwdctl-example.png)

To do so, you have to type this easy command: 

``` shell
sudo fwdctl create --destination-port 3000 --source-address 192.168.199.105 --source-port 80
```

That's it.

## Installation

#### Linux x86_64

```shell
wget https://github.com/alegrey91/fwdctl/releases/download/v0.1.0/fwdctl_Linux_x86_64 -O fwdctl && chmod +x fwdctl && sudo mv fwdctl /usr/local/bin/
```

