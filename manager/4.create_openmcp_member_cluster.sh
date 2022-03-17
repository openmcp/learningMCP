#!/bin/bash
if [ "$1" == "" ]; then
  echo "ClusterName Empty"
  exit 1
fi

for name in `kind get clusters | awk '{print $1}'`
do
  if [ $name == $1 ]; then
     echo "$1 is Already Exist"
     exit 1
  fi
done


flock -x /tmp/$1.lock -c "echo 'clusterCreating' > /home/keti/learningMCP/status/$1"


#[클러스터 생성]
chmod 666 /var/run/docker.sock
su $1 -c "kind create cluster --name $1 --config /home/keti/learningMCP/config/member-config.yaml"
#flock -x /tmp/kind.lock -c "su $1 -c 'kind create cluster --name $1 --config /home/keti/learningMCP/config/member-config.yaml'"
#su $1 -c "kind create cluster --name $1 --config /home/keti/learningMCP/config/member-config.yaml --loglevel=debug --retain"

su $1 -c "kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml"

#su $1 -c "kubectl create -f /home/keti/learningMCP/registry/cm_registry.yaml"
#/home/keti/learningMCP/registry/2.member_image_load.sh $1

/home/keti/learningMCP/console/run-butterfly.sh $1 &
flock -x /tmp/$1.lock -c "echo 'clusterCreated' > /home/keti/learningMCP/status/$1"

