
## 环境准备
yum install -y mariadb-server.x86_64

systemctl enable mariadb.service

systemctl start mariadb.service

cd /root

mkdir FinanceMgr

cd FinanceMgr


mysql < /root/FinanceMgr/src/analysis-server/api/sql/first_create_db.sql

mysql -e "GRANT ALL PRIVILEGES ON finance_mgr.* TO 'mgr'@'localhost' IDENTIFIED BY 'mgr'"

mysql -e "GRANT ALL PRIVILEGES ON finance_mgr.* TO 'mgr'@'%' IDENTIFIED BY 'mgr'"


## 代码编译及运行

export GOPATH=/root/

cd $GOPATH/FinanceMgr/src/analysis-server/api

go build -o analysis_server

./analysis_server

cd $GOPATH/FinanceMgr/src/analysis-server/cli

go build -o analysis_cli

souce openrc

