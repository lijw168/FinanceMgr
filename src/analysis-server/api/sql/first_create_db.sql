CREATE DATABASE IF NOT EXISTS `finance_mgr_2021` DEFAULT CHARACTER SET utf8;

/*==============================================================*/
/* Table: operatorInfo                                          */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`operatorInfo`;
create table if not exists `finance_mgr_2021`.`operatorInfo`
(
   `name`                 varchar(10) not null ,
   `password`             varchar(12),
   `companyId`            int not null,
   `job`                  varchar(32),
   `department`           varchar(64),
   `status`               int COMMENT '状态 ：0:offline;1:online;2:invalid user',
   `role`                 int COMMENT '角色 ：1：记账，2：审核，3：出纳，4：制单',
   `create_at`            datetime,
   `update_at`            datetime,
   primary key (name)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

-- use  finance_mgr_2021;
-- alter table operatorInfo add constraint FK_Reference_1 foreign key (companyId)  
--       references companyInfo (companyId) on delete restrict on update restrict;

/*==============================================================*/
/* 会计科目表 Table: accountSubject                              */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`accountSubject`;
create table if not exists `finance_mgr_2021`.`accountSubject`
(
   `subjectId`            int not null,
   `subjectName`          varchar(24) not null,
   `subjectLevel`         tinyint not null,
   primary key (subjectId),
   UNIQUE KEY `subjectName` (`subjectName`),
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* 凭证信息表 Table: voucherRecordInfo                               */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`voucherRecordInfo`;
create table  if not exists `finance_mgr_2021`.`voucherRecordInfo`
(
   `recordId`             int not null,
   `voucherId`            int  not null COMMENT '凭证ID',
   `subjectName`          varchar(64) not null COMMENT '会计科目名称，由1 ~ 4级的名称组合而成的',
   `debitMoney`           decimal(12,4) not null COMMENT '借方金额',
   `creditMoney`          decimal(12,4) not null COMMENT '贷方金额',
   `summary`              varchar(128) COMMENT '摘要',
   `subId1`               int COMMENT '一级会计科目ID',
   `subId2`               int COMMENT '二级会计科目ID',
   `subId3`               int COMMENT '三级会计科目ID',
   `subId4`               int COMMENT '四级会计科目ID',
   `billCount`            int COMMENT '该凭证记录的单据个数',
   `create_at`            datetime,
   `update_at`            datetime,
   primary key (recordId)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

-- alter table voucherRecordInfo add constraint FK_Reference_4 foreign key (voucherId)
--       references voucherInfo (voucherId) on delete restrict on update restrict;

/*==============================================================*/
/* 凭证信息表 Table: VoucherInfo                                 */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`voucherInfo`;
create table if not exists `finance_mgr_2021`.`voucherInfo`
(
   `voucherId`            int not null,
   `companyId`            int not null,   
   `voucherMonth`         int not null COMMENT '制证月份',
   `numOfMonth`           int not null COMMENT '本月第几次记录凭证',
   `voucherDate`          date not null COMMENT '制证日期',
   `create_at`            datetime,
   `update_at`            datetime,
   primary key (voucherId)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

-- alter table voucherInfo add constraint FK_Reference_3 foreign key (companyId)
--       references companyInfo (companyId) on delete restrict on update restrict;

/*==============================================================*/
/* Table: companyInfo                                           */
/*==============================================================*/
drop table if exists `finance_mgr_2021`.`companyInfo`;
create table if not exists `finance_mgr_2021`.`companyInfo`
(
   `companyId`            int not null,
   `companyName`          varchar(64),
   `abbreName`            varchar(24),
   `corporator`           varchar(16),
   `phone`                varchar(13),
   `e_mail`               varchar(32),
   `companyAddr`          varchar(128),
   `backup`               varchar(32),
   `create_at`            datetime,
   `update_at`            datetime,
   primary key (companyId),
   UNIQUE KEY `companyName` (`companyName`),
   UNIQUE KEY `abbreName` (`abbreName`)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* Table: IDInfo                                           */
/* companyId：从1开始，设计的值是到100； subjectId：从101开始，设计的值是到500；*/
/* voucherId：从501开始，设计的值的最大值，是int类型的最大值； */  
/* recordId：从1001开始。设计的值的最大值，是int类型的最大值；*/
/*==============================================================*/
create table if not exists `finance_mgr_2021`.`idInfo`
(
   `companyId`            int not null,
   `subjectId`            int not null,
   `voucherId`            int not null,
   `voucherRecordId`      int not null,
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;


