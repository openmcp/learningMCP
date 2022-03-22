#!/bin/bash
if [ "$1" == "" ]; then
  echo "ClusterName Empty"
  exit 1
fi

num=`echo $1 | sed 's/[^0-9]//g'`
port=575$num
#kill -9 $(ps -ef | grep python | grep $port | awk '{print $2}')
systemctl stop butterfly-$1.service
