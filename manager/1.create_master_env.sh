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

touch /tmp/$1.lock
chmod 666 /tmp/$1.lock

flock -x /tmp/$1.lock -c "touch /home/keti/learningMCP/status/$1"
chmod 666 /home/keti/learningMCP/status/$1
flock -x /tmp/$1.lock -c "echo 'envCreating' > /home/keti/learningMCP/status/$1"


#[유저 생성]
exec 3>/tmp/usrcreate.lock
flock -x 3

adduser --gecos "" --disabled-password $1
chpasswd <<< $1:1234

exec 3>&-

usermod -G keti $1
usermod -aG sudo $1

#[config 복사]
mkdir -p /home/$1/.kube
#cp -f /etc/kubernetes/admin.conf /home/$1/.kube/config
touch /home/$1/.kube/config
chown $1:$1 /home/$1/.kube
chown $1:$1 /home/$1/.kube/config
chmod go-r /home/$1/.kube/config

chown -R $1:$1 /home/$1/bin

su $1 -c "git clone https://github.com/openmcp/Public_OpenMCP-Release /home/$1/Public_OpenMCP-Release"
cd /home/keti/learningMCP/manager


sts_file="/home/keti/learningMCP/status/$1"

echo "exec 3>/tmp/$1.lock" >> /home/$1/.profile
echo "flock -x 3" >> /home/$1/.profile

echo "if [[ \`cat $sts_file\` == 'clusterCreated' ]]; then" >> /home/$1/.profile
echo "        echo 'connected' > $sts_file" >> /home/$1/.profile
echo "fi" >> /home/$1/.profile

echo "exec 3>&-" >> /home/$1/.profile

echo "exec 3>/tmp/$1.lock" >> /home/$1/.bash_logout
echo "flock -x 3" >> /home/$1/.bash_logout

echo "if [[ \`cat $sts_file\` == 'connected' ]]; then" >> /home/$1/.bash_logout
echo "        echo 'terminated' > $sts_file" >> /home/$1/.bash_logout
echo "fi" >> /home/$1/.bash_logout

echo "exec 3>&-" >> /home/$1/.bash_logout

echo "rm -rf /home/$1/Public_OpenMCP-Release/install_openmcp/master" >> /home/$1/.bash_logout
echo "rm -rf /home/$1/Public_OpenMCP-Release/install_openmcp/member" >> /home/$1/.bash_logout

#[User Specific hosts setting]
git clone https://github.com/figiel/hosts.git /home/$1/hosts
cd /home/$1/hosts
make
mkdir /home/$1/bin
cp libuserhosts.so /home/$1/bin
echo 'export LD_PRELOAD=~/bin/libuserhosts.so' >> /home/$1/.bashrc
touch /home/$1/.hosts
chown $1:$1 /home/$1/.hosts

flock -x /tmp/$1.lock -c "echo 'envCreated' > /home/keti/learningMCP/status/$1"

