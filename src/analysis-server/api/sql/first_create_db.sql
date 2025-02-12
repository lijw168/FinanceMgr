drop DATABASE if exists `finance_mgr`;
CREATE DATABASE IF NOT EXISTS `finance_mgr` DEFAULT CHARACTER SET utf8;


/*==============================================================*/
/* Table: companyGroup                                          */
/*==============================================================*/
drop table if exists `finance_mgr`.`companyGroup`;
create table if not exists `finance_mgr`.`companyGroup`
(
   `company_group_id`       int not null,
   `group_name`             varchar(64),
   `group_status`           int COMMENT '目前暂定两种状态：1：有效状态；0：无效状态',
   `created_at`             datetime,
   `updated_at`             datetime,
   primary key (company_group_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* Table: companyInfo                                           */
/*company_group_id,该字段为0的话，就表示该公司不与其他任何公司组成为组*/
/*==============================================================*/
drop table if exists `finance_mgr`.`companyInfo`;
create table if not exists `finance_mgr`.`companyInfo`
(
   `company_id`            int not null,
   `company_name`          varchar(64),
   `abbre_name`            varchar(24),
   `corporator`            varchar(16),
   `phone`                 varchar(13),
   `e_mail`                varchar(32),
   `company_addr`          varchar(128),
   `backup`                varchar(32),
   `start_account_period`  int not null COMMENT '启用会计期',
   `latest_account_year`   int not null COMMENT '最新的会计年度',
   `created_at`            datetime,
   `updated_at`            datetime,
   `company_group_id`      int DEFAULT 0 COMMENT '不能作为companyGroup的外键',
   primary key (company_id),
   UNIQUE KEY `company_name` (`company_name`),
   UNIQUE KEY `abbre_name` (`abbre_name`)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* Table: operatorInfo  */
/*role:第一个字节的低四位数值是操作员权限；第一个字节的高四位的数值是管理员的权限。*/
/*角色 ：1：制单, 2：出纳, 4：审核, 8：记账,16：查询, 32：增加，64：修改, 128:删除*/
/*==============================================================*/
drop table if exists `finance_mgr`.`operatorInfo`;
create table if not exists `finance_mgr`.`operatorInfo`
(
   `operator_id`          int not null,
   `name`                 varchar(10) not null COMMENT '一个公司内，不允许有重复的',
   `password`             varchar(12),
   `company_id`           int not null,
   `job`                  varchar(32),
   `department`           varchar(64),
   `status`               int DEFAULT 0 COMMENT '状态 ：0:invalid status;1:online;2:offline;',
   `role`                 int DEFAULT 0 COMMENT '角色',
   `created_at`           datetime,
   `updated_at`           datetime,
   primary key (operator_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

create index comId_name on `finance_mgr`.`operatorInfo` (company_id,name);

use  finance_mgr;
alter table operatorInfo add constraint FK_Reference_1 foreign key (company_id)  
      references companyInfo (company_id) on delete restrict  on update restrict;

/*==============================================================*/
/* 会计科目表 Table: accountSubject                              */
/* common_id是在一个公司内，不允许有重复的;但subject_name会计科目的名称，*/
/* 在一个公司内，1级科目名称不能重复，但二级以后的科目名称是可以重复的。*/
/*subject_type:科目的类型;0:不选择；1:资产;2:负债;3:权益;4:成本;5:损益*/
/*subject_direction:科目的性质;1:debit;2:credit*/
/*==============================================================*/
drop table if exists `finance_mgr`.`accountSubject`;
create table if not exists `finance_mgr`.`accountSubject`
(
   `subject_id`            int not null,
   `company_id`            int not null,
   `common_id`             varchar(10) not null COMMENT '该ID是操作用户添加的，该行业习惯用的ID' ,
   `subject_name`          varchar(24) not null ,
   `subject_level`         tinyint not null,
   `subject_direction`     tinyint not null,
   `subject_type`          tinyint not null,
   `mnemonic_code`         varchar(10) not null,
   `subject_style`         varchar(10) not null,
   primary key (subject_id)
   /*unique key `subjectName` (`subject_name`),*/
   /*unique key `commonId` (`common_id`)*/
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

alter table accountSubject add constraint FK_Reference_2 foreign key (company_id)
      references companyInfo (company_id) on delete restrict  on update restrict;

/*==============================================================*/
/* Table: 该表保存的数据是各个科目的年初余额，一般只有资产类、负债和权益类才有年初余额，损益和成本没有。
   status : 1，已结算；0，未结算
   balance： 正数，表示该与该科目本身的方向一致。2、负数，表示该值与该科目本身的方向相反*/
/*==============================================================*/
drop table if exists `finance_mgr`.`beginOfYearBalance`;
create table if not exists `finance_mgr`.`beginOfYearBalance`
(
   `id`                    int not null AUTO_INCREMENT,
   `company_id`            int not null,
   `subject_id`            int not null,
   `year`                  int not null,
   `balance`               decimal(12,4) not null,
   `status`                tinyint not null default 0,
   `created_at`            datetime,
   `updated_at`            datetime,
   primary key (id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

create index comId_subId_year_index on `finance_mgr`.`beginOfYearBalance` (company_id,year,subject_id);

-- 之所以删除该约束，就是因为这两个表之间没有强关联性，所以删除
-- alter table beginOfYearBalance add constraint FK_Reference_5 foreign key (subjectId)
--       references accountSubject (subjectId) on delete restrict;

/*==============================================================*/
/* 凭证信息表 Table: VoucherInfo                                 */
/*==============================================================*/
drop table if exists `finance_mgr`.`voucherInfo`;
create table if not exists `finance_mgr`.`voucherInfo`
(
   `voucher_id`            int not null,
   `company_id`            int not null, 
   `voucher_month`         int not null COMMENT '制证月份',
   `num_of_month`          int not null COMMENT '本月第几次记录凭证',
   `voucher_date`          int not null COMMENT '制证日期',
   `voucher_filler`        varchar(10) COMMENT '制证者',
   `voucher_auditor`       varchar(10) COMMENT '审核者',
   `bill_count`            int DEFAULT 0 COMMENT '该张凭证的单据个数',
   `status`                int DEFAULT 1 COMMENT '1:未审核；2：已作废；3：已审核 ...',
   `created_at`            datetime,
   `updated_at`            datetime,
   primary key (voucher_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

alter table voucherInfo add constraint FK_Reference_3 foreign key (company_id)
      references companyInfo (company_id) on delete restrict  on update restrict;

/*==============================================================*/
/* 凭证信息表 Table: voucherRecordInfo                          */
/* 目前该表中只使用了sub_id1 表示科目的ID。 其他三个ID字段没有使用 */
/*==============================================================*/
drop table if exists `finance_mgr`.`voucherRecordInfo`;
create table  if not exists `finance_mgr`.`voucherRecordInfo`
(
   `record_id`             int not null,
   `voucher_id`            int  not null COMMENT '凭证ID',
   `subject_name`          varchar(128)  not null COMMENT '会计科目名称，由1 ~ 4级的名称组合而成的',
   `debit_money`           decimal(12,4) not null COMMENT '借方金额',
   `credit_money`          decimal(12,4) not null COMMENT '贷方金额',
   `summary`               varchar(128) COMMENT '摘要',
   `sub_id1`               int DEFAULT 0 COMMENT '一级会计科目ID',
   `sub_id2`               int DEFAULT 0 COMMENT '二级会计科目ID',
   `sub_id3`               int DEFAULT 0 COMMENT '三级会计科目ID',
   `sub_id4`               int DEFAULT 0 COMMENT '四级会计科目ID',
   `created_at`            datetime,
   `updated_at`            datetime,
   primary key (record_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

alter table voucherRecordInfo add constraint FK_Reference_4 foreign key (voucher_id)
      references voucherInfo (voucher_id) on delete restrict on update restrict;


/*==============================================================*/
/* Table: voucherTemplate                                */
/* 之所以增加company_id,为了方便获取某个公司的voucherTemplate*/
/* 可以增加一个非唯一索引，在company_id字段上，因为list时，是通过company_id*/
/*==============================================================*/
drop table if exists `finance_mgr`.`voucherTemplate`; 
create table if not exists `finance_mgr`.`voucherTemplate`
(
   `voucher_template_id`    int not null,
   `company_id`             int not null,
   `reference_voucher_id`   int not null COMMENT '所引用的凭证ID',
   `voucher_year`           int not null COMMENT '所引用的凭证数据年度',
   `illustration`           varchar(24),
   `created_at`             datetime,
   primary key (voucher_template_id),
   index comId_Year_name (company_id,voucher_year)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

create index comId_index on `finance_mgr`.`voucherTemplate` (company_id);

/*==============================================================*/
/* Table: IDInfo                                           */
/* companyId：从1开始，设计的值是到100；operator_id:从101开始*/
/* subjectId：从501开始，设计的值是到1000；*/
/* voucherId：从1001开始，设计的值的最大值，是int类型的最大值； */  
/* recordId：从5001开始。设计的值的最大值，是int类型的最大值；*/
/* companyGroupId：从801开始。设计的值的最大值，是int类型的最大值；*/
/* voucherTemplateId：从1开始。设计的值的最大值，是int类型的最大值；*/
/*==============================================================*/
drop table if exists `finance_mgr`.`idInfo`;
create table if not exists `finance_mgr`.`idInfo`
(
   `company_id`             int not null,
   `operator_id`            int not null,
   `subject_id`             int not null,
   `voucher_id`             int not null,
   `voucher_record_id`      int not null,
   `company_group_id`       int not null,
   `voucher_template_id`    int not null
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* Table: userLoginInfo                                         */
/*==============================================================*/
drop table if exists `finance_mgr`.`userLoginInfo`;
create table if not exists `finance_mgr`.`userLoginInfo`
(
   `id`                   int primary key AUTO_INCREMENT,
   `operator_id`          int not null,
   `name`                 varchar(10) not null ,
   `status`               int DEFAULT 0 COMMENT '状态 ：0:invalid status;1:online;2:offline;',
   `client_ip`            varchar(16),
   `begined_at`           datetime,
   `ended_at`             datetime
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

/*==============================================================*/
/* Table: menuInfo                                              */   
/*menu_serial_num:用于排列菜单的顺序。*/                                           
/*==============================================================*/
drop table if exists `finance_mgr`.`menuInfo`;
create table if not exists `finance_mgr`.`menuInfo`
(
   `menu_id`              int not null,
   `menu_name`            varchar(24),
   `menu_level`           int,
   `parent_menu_id`       int,
   `menu_serial_num`      int,
   primary key (menu_id)
)ENGINE=InnoDB DEFAULT CHARSET=UTF8;

insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(71,"系统",1,0,1);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(7101,"系统设置",2,71,1000);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(7102,"菜单管理",2,71,1001);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(7103,"重新登录",2,71,1002);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(7104,"修改密码",2,71,1003);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(7105,"退出",2,71,1004);

insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(81,"总账系统",1,0,2);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8101,"明细账",2,81,1101);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8102,"余额表",2,81,1102);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8103,"科目汇总",2,81,1103);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8104,"日记账",2,81,1104);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8105,"总账",2,81,1105);

insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8120,"记账凭证",2,81,1120);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8122,"审核凭证",2,81,1121);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8121,"会计科目",2,81,1122);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8124,"年度结算",2,81,1123);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8125,"取消年度结算",2,81,1124);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(8123,"年初余额",2,81,1125);

insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(51,"项目管理",1,0,3);

insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(61,"工具",1,0,4);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(6101,"word",2,61,1200);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(6102,"excel",2,61,1211);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(6103,"计算器",2,61,1212);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(6104,"记事本",2,61,1213);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(41,"帮助",1,0,5);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(4101,"在线升级",2,41,1220);
insert into menuInfo(menu_id,menu_name,menu_level,parent_menu_id,menu_serial_num) value(4102,"关于...",2,41,1221);


insert into companyInfo(company_id,company_name,abbre_name,corporator,phone,e_mail,company_addr,backup,created_at,updated_at) value(1,"rootManager","manager","","","","","",now(),now());
insert into operatorInfo (operator_id,name,password,company_id,job,department,status,role,created_at,updated_at) value(101,"root","root@123",1,"maintainer","",1,255,now(),now());
insert into idInfo (company_id,operator_id,subject_id,voucher_id,voucher_record_id,company_group_id,voucher_template_id) value(2,102,501,1001,5001,801,1);
