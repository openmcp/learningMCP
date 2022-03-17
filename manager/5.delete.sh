#!/bin/bash
if [ "$1" == "" ]; then
  echo "ClusterName Empty"
  exit 1
fi

if [ "$1" == "keti" ]; then
  echo "can't not remove user: keti"
  exit 1
fi
/home/keti/learningMCP/console/stop-butterfly.sh $1

flock -x /tmp/$1.lock -c "echo 'deleting' > /home/keti/learningMCP/status/$1"

kind delete cluster --name $1
cat /dev/null > /home/$1/.kube/config
flock -x /tmp/$1.lock -c "echo 'deleted' > /home/keti/learningMCP/status/$1"
