#!/bin/bash
if [ "$1" == "" ]; then
  echo "ClusterName Empty"
  exit 1
fi

if [ "$1" == "keti" ]; then
  echo "can't not remove user: keti"
  exit 1
fi

flock -x /tmp/$1.lock -c "echo 'deleting' > /home/keti/learningMCP/status/$1"

/home/keti/learningMCP/console/stop-butterfly.sh $1

kind delete cluster --name $1
userdel -rfRZ $1
rm -rf /home/$1
flock -x /tmp/$1.lock -c "rm /home/keti/learningMCP/status/$1"
groupdel $1
