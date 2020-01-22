#!/usr/bin/env bash
echo ON
export GOPATH=$(dirname "$PWD"):/Users/yepeng/go

./autoVersionTool_mac ../src/version/ver.go
cd ~/project/resourceServer/src

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o resourceServer_new

if [ ${1} -eq 205 ];then
    scp ~/project/resourceServer/src/resourceServer_new root@192.168.3.205:/root/servers/resourceServer
    ssh root@192.168.3.205 "cd /root/servers/resourceServer;mv resourceServer_new resourceServer; ./start.sh"
    echo "update 205 success"
elif [ ${1} -eq 220 ];then
    scp ~/project/resourceServer/src/resourceServer_new root@192.168.1.220:/opt/resourceServer
    ssh root@192.168.1.220 "cd /opt/resourceServer; mv resourceServer_new resourceServer; ./start.sh"
    echo "update 220 success"
else
   echo "error"
fi

