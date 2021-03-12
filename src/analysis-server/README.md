
## 环境准备
yum install -y mariadb-server.x86_64

systemctl enable mariadb.service

systemctl start mariadb.service

cd /root

mkdir zbs

cd zbs

git clone git@git.jd.com:iaas-sdn/zbs-server.git

mysql < /root/zbs/zbs-server/api/sql/first_create_db.sql

mysql -e "GRANT ALL PRIVILEGES ON finance_mgr.* TO 'lijw'@'localhost' IDENTIFIED BY 'lijw'"

mysql -e "GRANT ALL PRIVILEGES ON finance_mgr.* TO 'lijw'@'%' IDENTIFIED BY 'lijw'"


## 代码编译及运行

export GOPATH=/root/

cd $GOPATH/zbs/src/analysis-server/api

go build -o zbs-server

./zbs-server

cd $GOPATH/jcloud-zbs/src/analysis-server/cli

go build -o zbs_cli

souce openrc

./zbs_cli pool-create pool1 ssd 16 104857600

./zbs_cli pool-list

./zbs_cli rack-create tag1

./zbs_cli rack-list

./zbs_cli host-create vrouter1 rack-dd2iux8odk pool-m95mr7hgt6 --manage_ip 10.12.209.161 --storage_ip 10.12.209.161 --client_ip 10.12.209.161 --trace_ip 10.12.209.161

./zbs_cli host-list

./zbs_cli disk-create 24097d56-2df8-4e9c-ba4f-ec3f3aa0e52c host-2su8leg40r ssd 10737418240 --manage_ip 10.12.209.161 --storage_ip 10.12.209.161 --client_ip 10.12.209.161 --trace_ip 10.12.209.161

./zbs_cli disk-list

./zbs_cli MoveReplica replicaId diskId

./zbs_cli Rescheduler poolId

./zbs_cli TransferLeader ReplicationGroupId ReplicaId

# TODO
+ Dao 里面的transaction 错误必须处理
+ 给Disk启动的时候提供ObjectSize
+ VolumeAPI response 返回增加详细信息
+ ZBS->EBS 恢复Snapshot查询状态不更新
+ describeVolume需要attachment
