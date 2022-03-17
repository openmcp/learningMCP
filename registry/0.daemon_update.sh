#vim /etc/docker/daemon.json 
# insert => "insecure-registries": ["10.0.3.40:5000"]

systemctl daemon-reload
systemctl restart docker

