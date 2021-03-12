#!/bin/bash

set -e

origin_db=zbs
global_db=zbs_global 
az_db=zbs_az

echo "WARNING: this will destroy zbs tables in database zbs_global and zbs_az"

doit=0

if [ "$#" != "1" ]; then
	echo "$0 [--dryrun | --doit]"
	exit 1
elif [ "$1" != "--dryrun" ] && [ "$1" != "--doit" ];then
	echo "$0 [--dryrun | --doit]"
	exit 1
fi

if [ "$1" == "--doit" ];then
	doit=1
fi

global_tables=(
quota             \
snapshot          \
snapshot_task     \
tasks             \
volume            \
volume_attachment \
volume_task       \
volume_type
)

az_tables=(
admin_tasks       \
disk              \
host              \
pool              \
proxy             \
rack              \
replica           \
replica_actions   \
replica_status    \
replication_group \
replication_task  \
sched_elect       \
tenant_pool_map
)

function copy_table() {
	if [ "$#" != "4" ];then
		echo "copy_table function should have four arguments"
		exit 1
	fi
	db_src=$1
	tb_src=$2
	db_dest=$3
	tb_dest=$4
	if [ "$doit" == "1" ];then
		mysql -e "DROP TABLE IF EXISTS $db_dest.$tb_dest"
		mysql -e "CREATE TABLE $db_dest.$tb_dest LIKE $db_src.$tb_src"
		mysql -e "INSERT INTO $db_dest.$tb_dest SELECT * FROM $db_src.$tb_src"
	else
		echo "DROP TABLE IF EXISTS $db_dest.$tb_dest"
		echo "CREATE TABLE $db_dest.$tb_dest LIKE $db_src.$tb_src"
		echo "INSERT INTO $db_dest.$tb_dest SELECT * FROM $db_src.$tb_src"
	fi
}

if [ "$doit" == "1" ];then
	mysql -e "CREATE DATABASE IF NOT EXISTS $global_db"
	mysql -e "CREATE DATABASE IF NOT EXISTS $az_db"
	mysql -e "GRANT ALL PRIVILEGES ON zbs_global.* TO 'zbs_global'@'127.0.0.1' IDENTIFIED BY 'zbs_global'"
	mysql -e "GRANT ALL PRIVILEGES ON zbs_global.* TO 'zbs_global'@'%' IDENTIFIED BY 'zbs_global'"
	mysql -e "GRANT ALL PRIVILEGES ON zbs_az.* TO 'zbs_az'@'127.0.0.1' IDENTIFIED BY 'zbs_az'"
	mysql -e "GRANT ALL PRIVILEGES ON zbs_az.* TO 'zbs_az'@'%' IDENTIFIED BY 'zbs_az'"
else
	echo "CREATE DATABASE IF NOT EXISTS $global_db"
	echo "CREATE DATABASE IF NOT EXISTS $az_db"
fi

for i in ${global_tables[@]}; do 
	copy_table $origin_db $i $global_db $i
done

for i in ${az_tables[@]}; do 
	copy_table $origin_db $i $az_db $i
done
