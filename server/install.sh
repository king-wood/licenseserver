#! /bin/bash
set -e -x
docker stop $(docker ps -qa)
docker rm $(docker ps -qa)
docker rmi licenseserver:latest

docker load < licenseserver.tar

docker run -d --name licenseserver -v `pwd`/log/:/log/ -v `pwd`/license.db:/license.db -p 8080:8080 licenseserver
