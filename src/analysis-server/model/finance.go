package model

import (
	"time"
)

//OperatorInfo ...
type OperatorInfo struct {
	CompanyID  int    `orm:"column(companyId)"`
	Name       string `orm:"pk;column(name)"`
	Password   string `orm:"column(password)"`
	Job        string `orm:"column(job)"`
	Department string `orm:"column(department)"`
	Status     int    `orm:"column(Status)"`
	Role       int    `orm:"column(role)"`
}

//VoucherRecord ...
type VoucherRecord struct {
	RecordID    int     `orm:"pk;column(recordId)"`
	VoucherID   int     `orm:"column(voucherId)"`
	SubjectName string  `orm:"column(subjectName)"`
	DebitMoney  float64 `orm:"column(debitMoney)"`
	CreditMoney float64 `orm:"column(creditMoney)"`
	Summary     string  `orm:"column(summary)"`
	SubID1      int     `orm:"column(subId1)"`
	SubID2      int     `orm:"column(subId2)"`
	SubID3      int     `orm:"column(subId3)"`
	SubID4      int     `orm:"column(subId4)"`
	BillCount   int     `orm:"column(billCount)"`
}

//AccSubject ...
type AccSubject struct {
	SubjectID    string `orm:"pk;column(subjectId)"`
	SubjectName  string `orm:"column(subjectName)"`
	SubjectLevel int    `orm:"column(subjectLevel)"`
}

//VoucherInfo ...
type VoucherInfo struct {
	VoucherID    int       `orm:"pk;column(voucherId)"`
	CompanyID    int       `orm:"column(companyId)"`
	VoucherMonth int       `orm:"column(voucherMonth)"`
	NumOfMonth   int       `orm:"column(numOfMonth)"`
	VoucherDate  time.Time `orm:"column(voucherDate)"`
}

//CompanyInfo ...
type CompanyInfo struct {
	CompanyID   int    `orm:"pk;column(companyId)"`
	CompanyName string `orm:"column(companyName)"`
	AbbrevName  string `orm:"column(abbreviationName)"`
	Corporator  string `orm:"column(corporator)"`
	Phone       string `orm:"column(phone)"`
	Summary     string `orm:"column(summary)"`
	Email       string `orm:"column(e_mail)"`
	CompanyAddr string `orm:"column(companyAddr)"`
	Backup      string `orm:"column(backup)"`
}

//IDInfo ...
type IDInfo struct {
	VoucherID       int `orm:"column(voucherId)"`
	CompanyID       int `orm:"column(companyId)"`
	SubjectID       int `orm:"column(subjectId)"`
	VoucherRecordID int `orm:"column(voucherRecordId)"`
}

//CommResp ...
type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}
