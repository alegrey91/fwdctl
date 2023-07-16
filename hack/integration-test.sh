#!/bin/sh
#
set -e

# create directory to store cover profiles
INTEGDIR=/tmp/integration
mkdir -p $INTEGDIR
export GOCOVERDIR=$INTEGDIR

# list of subcommands we want to test
command_list="version|list|help|generate rules -O /tmp/|generate systemd -O /tmp/ -f /tmp/rules.yaml -p /tmp/ -t fork|add -d 8080 -i lo -P tcp -s 127.0.0.1 -p 9090|delete -n 1"

IFS="|"
for cmd in $command_list;
do
  command="./fwdctl $cmd"
  echo "$command"
  eval "$command"
done
