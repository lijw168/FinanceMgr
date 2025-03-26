
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

go build -o analysis-server

./analysis_server

cd $GOPATH/FinanceMgr/src/analysis-server/cli

go build -o analysis_cli

souce openrc

./analysis_cli ListAccSub
## 要考虑一下ID达到最大值后，怎样进行处理；需要修改初始化的ID值时，怎样进行修改？
## 目前先通过手动插入一下那几个ID值的最小值。

## 删除公司的数据，所涉及到的表
select start_account_period,latest_account_year from companyInfo where company_id = 7;

select * from beginOfYearBalance where company_id = 7;
delete from beginOfYearBalance where company_id = 7;

select * from voucherInfo_2022 where company_id = 7;
delete from voucherInfo_2022 where company_id = 7;

select * from voucherRecordInfo_2022 where voucher_id in (select voucher_id from voucherInfo_2022 where company_id = 7);
delete from voucherRecordInfo_2022 where voucher_id in (select voucher_id from voucherInfo_2022 where company_id = 7);

select * from voucherTemplate where company_id = 7;

