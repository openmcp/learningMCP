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
su $1 -c "kind create cluster --name $1 --config /home/keti/learningMCP/config/master-config-$1.yaml"
#flock -x /tmp/kind.lock -c "su $1 -c 'kind create cluster --name $1 --config /home/keti/learningMCP/config/master-config-$1.yaml'"
#su $1 -c "kind create cluster --name $1 --config /home/keti/learningMCP/config/master-config-$1.yaml --loglevel=debug --retain"

su $1 -c "kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml"

chown $1:$1 /home/$1/.init

cd /home/keti/learningMCP/create_cluster
su $1 -c "./get_helm.sh"
su $1 -c "./6.install_federation.sh"
cd /home/keti/learningMCP/manager

su $1 -c "kubectl delete validatingwebhookconfiguration validations.core.kubefed.io"
#/home/keti/learningMCP/registry/1.master_image_load.sh $1

/home/keti/learningMCP/console/run-butterfly.sh $1 &
flock -x /tmp/$1.lock -c "echo 'clusterCreated' > /home/keti/learningMCP/status/$1"

