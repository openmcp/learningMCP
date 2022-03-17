[Install LearningMCP Env]
1. remove status
  -> rm /home/keti/learningMCP/status/*

2. docker registy connection
  -> cd /home/keti/learningMCP/registry
  -> read 0.daemon_update.sh && run
  -> run ./1.docker_login.sh
  -> run ./2.create_registry.sh

3. external Pdns Install (with docker)
  -> cd /home/keti/learningMCP/pdns
  -> ./pdns_install.sh

4. run daemon
  -> cd /home/keti/learningMCP/manager
  -> go run daemon.go

5. run butterfly
  -> cd /home/keti/learningMCP/console
  -> ./run-butterlfy.sh

[Maintance]
1. github source Pull (Only OpenMCP Master)
  -> cd ~/PublicOpenMCP
  -> git pull origin master

2. local registry
  -> cd /home/keti/learningMCP/registry
  -> ./3.image_push.sh



