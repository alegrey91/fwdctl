# fwdctl

<p align="center">
    <img src="fwdctl.png" alt="fwdctl" width="200"/>
</p>


[![Go Reference](https://pkg.go.dev/badge/github.com/alegrey91/fwdctl.svg)](https://pkg.go.dev/github.com/alegrey91/fwdctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/alegrey91/fwdctl)](https://goreportcard.com/report/github.com/alegrey91/fwdctl)
[![Go Coverage](https://github.com/alegrey91/fwdctl/wiki/coverage.svg)](https://raw.githack.com/wiki/alegrey91/fwdctl/coverage.html)
![https://img.shields.io/ossf-scorecard/github.com/alegrey91/fwdctl?label=openssf%20scorecard&style=flat](https://img.shields.io/ossf-scorecard/github.com/alegrey91/fwdctl?label=openssf%20scorecard&style=flat)
[![Awesome](https://awesome.re/badge.svg)](https://github.com/avelino/awesome-go/)

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

Full **documentation** [here](docs/getting-started.md).

## Installation

#### Linux x86_64

```shell
curl -s https://raw.githubusercontent.com/alegrey91/fwdctl/main/install | sudo sh
```

## Seccomp (experimental)

I've recently added a new functionality to trace the system calls used by `fwdctl` during the test pipeline.

This is done by using another project of mine: [`harpoon`](https://github.com/alegrey91/harpoon).

Thanks to this, at the end of the pipeline, we have a **seccomp** profile as artifact. You can use this to run `fwdctl` in a more secure way.

Find the **seccomp** profile here: [`fwdctl-seccomp.json`](hack/fwdctl-seccomp.json).
