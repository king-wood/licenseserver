#!/bin/bash
set -e -x
rm -fr licenseserver.tar
docker rmi licenseserver:latest

docker build -t licenseserver .
docker save licenseserver > licenseserver.tar

rsync -P -rz licenseserver.tar licenseserver:/home/ubuntu/licenseserver.tar

sqlite3 license.db < db/db.sql
rsync -P -rz license.db licenseserver:/home/ubuntu/license.db
