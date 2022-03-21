#!/bin/bash
HOST_IP=`ip route get 8.8.8.8 | head -1 | cut -d' ' -f8`   # input target server IP (ex.10.0.0.226)
HOST_NAME=learningmcp  # input target Server Hostname (ex.nanmdev6)
ETCD_VER=v3.4.3   #kube 1.17에서 사용하는 버전

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

mv /tmp/etcd-download-test ./etcd-3.4.3
cd ./etcd-3.4.3

echo "install success"

## etcd 실행
nohup ./etcd  --advertise-client-urls=http://${HOST_IP}:12379 --initial-advertise-peer-urls=http://${HOST_IP}:12380 --initial-cluster=${HOST_NAME}=http://${HOST_IP}:12380 --listen-client-urls=http://127.0.0.1:12379,http://${HOST_IP}:12379 --listen-metrics-urls=http://127.0.0.1:12381 --listen-peer-urls=http://${HOST_IP}:12380 --name ${HOST_NAME} &

echo "running success"
## 확인
#cat nohup
