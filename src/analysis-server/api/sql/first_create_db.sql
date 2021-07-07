CREATE DATABASE IF NOT EXISTS `finance_mgr_2021` DEFAULT CHARACTER SET utf8;

/*==============================================================*/
/* Table: companyInfo                                           */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`companyInfo`;
create table if not exists `finance_mgr_2021`.`companyInfo`
(
   `company_id`            int not null,
   `company_name`          varchar(64),
   `abbre_name`            varchar(24),
   `corporator`            varchar(16),
   `phone`                 varchar(13),
   `e_mail`                varchar(32),
   `company_addr`          varchar(128),
   `backup`                varchar(32),
   `created_at`            datetime,
   `updated_at`            datetime,
   primary key (company_id),
   UNIQUE KEY `company_name` (`company_name`),
   UNIQUE KEY `abbre_name` (`abbre_name`)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* Table: operatorInfo                                          */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`operatorInfo`;
create table if not exists `finance_mgr_2021`.`operatorInfo`
(
   `name`                 varchar(10) not null ,
   `password`             varchar(12),
   `company_id`           int not null,
   `job`                  varchar(32),
   `department`           varchar(64),
   `status`               int COMMENT '状态 ：0:invalid status;1:online;2:offline;',
   `role`                 int COMMENT '角色 ：1：记账，2：审核，3：出纳，4：制单',
   `created_at`           datetime,
   `updated_at`           datetime,
   primary key (name)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

use  finance_mgr_2021;
alter table operatorInfo add constraint FK_Reference_1 foreign key (company_id)  
      references companyInfo (company_id) on delete restrict on update restrict;

/*==============================================================*/
/* 会计科目表 Table: accountSubject                              */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`accountSubject`;
create table if not exists `finance_mgr_2021`.`accountSubject`
(
   `subject_id`            int not null,
   `common_id`             varchar(10) not null,   '该ID是操作用户添加的，该行业习惯用的ID'
   `subject_name`          varchar(24) not null,
   `subject_level`         tinyint not null,
   primary key (subject_id),
   unique key `subjectName` (`subject_name`)
   unique key `commonId` (`common_id`);
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* 凭证信息表 Table: VoucherInfo                                 */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`voucherInfo`;
create table if not exists `finance_mgr_2021`.`voucherInfo`
(
   `voucher_id`            int not null,
   `company_id`            int not null,   
   `voucher_month`         int not null COMMENT '制证月份',
   `num_of_month`          int not null COMMENT '本月第几次记录凭证',
   `voucher_date`          date not null COMMENT '制证日期',
   `created_at`            datetime,
   `updated_at`            datetime,
   primary key (voucher_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

alter table voucherInfo add constraint FK_Reference_3 foreign key (company_id)
      references companyInfo (company_id) on delete restrict on update restrict;

/*==============================================================*/
/* 凭证信息表 Table: voucherRecordInfo                               */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`voucherRecordInfo`;
create table  if not exists `finance_mgr_2021`.`voucherRecordInfo`
(
   `record_id`             int not null,
   `voucher_id`            int  not null COMMENT '凭证ID',
   `subject_name`          varchar(64) not null COMMENT '会计科目名称，由1 ~ 4级的名称组合而成的',
   `debit_money`           decimal(12,4) not null COMMENT '借方金额',
   `credit_money`          decimal(12,4) not null COMMENT '贷方金额',
   `summary`               varchar(128) COMMENT '摘要',
   `sub_id1`               int COMMENT '一级会计科目ID',
   `sub_id2`               int COMMENT '二级会计科目ID',
   `sub_id3`               int COMMENT '三级会计科目ID',
   `sub_id4`               int COMMENT '四级会计科目ID',
   `bill_count`            int COMMENT '该凭证记录的单据个数',
   `created_at`            datetime,
   `updated_at`            datetime,
   primary key (record_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

alter table voucherRecordInfo add constraint FK_Reference_4 foreign key (voucher_id)
      references voucherInfo (voucher_id) on delete restrict on update restrict;

/*==============================================================*/
/* Table: IDInfo                                           */
/* companyId：从1开始，设计的值是到100； subjectId：从101开始，设计的值是到500；*/
/* voucherId：从501开始，设计的值的最大值，是int类型的最大值； */  
/* recordId：从1001开始。设计的值的最大值，是int类型的最大值；*/
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`idInfo`;
create table if not exists `finance_mgr_2021`.`idInfo`
(
   `company_id`             int not null,
   `subject_id`             int not null,
   `voucher_id`             int not null,
   `voucher_record_id`      int not null
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

--insert into idInfo (company_id,subject_id,voucher_id,voucher_record_id) value(100,101,501,1001);

/*==============================================================*/
/* Table: userLoginInfo                                         */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`userLoginInfo`;
create table if not exists `finance_mgr_2021`.`userLoginInfo`
(
   `id`                   int primary key AUTO_INCREMENT,
   `name`                 varchar(10) not null ,
   `status`               int COMMENT '状态 ：1:offline;2:online;3:invalid user',
   `client_ip`            varchar(16),
   `begined_at`           datetime,
   `ended_at`             datetime
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;


