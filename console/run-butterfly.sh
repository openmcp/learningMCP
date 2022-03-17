#!/bin/bash
#myip=`ip route get 8.8.8.8 | head -1 | cut -d' ' -f8`
#docker run --name butterfly -d -p 57575:57575 garland/butterfly --unsecure --port=57575 #--login

# python 3.7 install
# pip3.7 install
# pip3 install butterfly

# https://github.com/paradoxxxzero/butterfly
#butterfly.server.py --unsecure --host=0.0.0.0 --port=57575 --login

#cp /home/keti/learningMCP/console/butterfly.service /etc/systemd/system
#chmod 777 /etc/systemd/system/butterfly.service
#systemctl daemon-reload
#systemctl enable butterfly.service
#systemctl start butterfly.service

if [ "$1" == "" ]; then
  echo "ClusterName Empty"
  exit 1
fi

num=`echo $1 | sed 's/[^0-9]//g'`
su $1 -c "/home/keti/learningMCP/console/butterfly/butterfly.server.py --unsecure --host=0.0.0.0 --port=575$num --debug"

