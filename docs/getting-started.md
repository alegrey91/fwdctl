# Documentation

**fwdctl** has many commands to manage your forwards. Let's dive into them!

### Apply

The `apply` command is used to apply rules from a rules file. Here's a rules file example.

```shell
cat << EOF > rules.yml
rules:
  - dport: 2022
    saddr: 192.168.122.43
    sport: 22
    iface: eth1
    proto: TCP
  - dport: 3022
    saddr: 192.168.122.44
    sport: 22
    iface: eth1
    proto: tcp
EOF
```

Once created, you can easily run:

```shell
sudo fwdctl apply --rules-file rules.yml
```

and apply all the rules listed in the file.

### Create

The `create` command is used to create single rules manually. Let's see an example.

![](../fwdctl-example.png)

To implement the rule for this scenario, just type the following command:

```shell
sudo fwdctl create --destination-port 3000 --source-address 192.168.199.105 --source-port 80
```

### Delete

The `delete` command is used to delete single rules manually.

To delete a specific rule (identified with a number), type the following command:

```shell
sudo fwdctl delete -n 4
```

This will remove the rule n. **4**, listed from **iptables**.

### List

The `list` command is used to list the applied rules.

To list the rules, type the following command:

```shell
sudo fwdctl list
```

The output will look like this:

```shell
+--------+-----------+----------+---------------+----------------+---------------+
| NUMBER | INTERFACE | PROTOCOL | EXTERNAL PORT |  INTERNAL IP   | INTERNAL PORT |
+--------+-----------+----------+---------------+----------------+---------------+
|      1 | lo        | tcp      |          2022 | 192.168.122.43 |            22 |
|      2 | lo        | tcp      |          3022 | 192.168.122.44 |           443 |
+--------+-----------+----------+---------------+----------------+---------------+
```

### Generate

The `generate` command is used to generate the following files:

* **systemd** service for `fwdctl`
* rules empty file

To generate a systemd service, type the following command:

```shell
fwdctl generate systemd -o fwdctl.service
```

To generate a `fwdctl` rules file, instead, type the following command:

```shell
fwdctl generate rules -o rules.yml
```

### Version

The `version` command is used show the version of the program.
