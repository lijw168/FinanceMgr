package model

import (
	"time"
)

//LoginInfo ...
type LoginInfo struct {
	OperatorID int
	Name       string
	Status     int
	ClientIp   string
	BeginedAt  time.Time
	EndedAt    time.Time
}

//OperatorInfo ...
type OperatorInfo struct {
	OperatorID int
	CompanyID  int
	Name       string
	Password   string
	Job        string
	Department string
	Status     int
	Role       int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

//VoucherRecord ...
type VoucherRecord struct {
	RecordID    int
	VoucherID   int
	SubjectName string
	DebitMoney  float64
	CreditMoney float64
	Summary     string
	SubID1      int
	SubID2      int
	SubID3      int
	SubID4      int
	BillCount   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//AccSubject ...
type AccSubject struct {
	SubjectID    int
	CompanyID    int
	CommonID     string
	SubjectName  string
	SubjectLevel int
}

//VoucherInfo ...
type VoucherInfo struct {
	VoucherID    int
	CompanyID    int
	VoucherMonth int
	NumOfMonth   int
	VoucherDate  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//CompanyInfo ...
type CompanyInfo struct {
	CompanyID   int
	CompanyName string
	AbbrevName  string
	Corporator  string
	Phone       string
	Email       string
	CompanyAddr string
	Backup      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//IDInfo ...
type IDInfo struct {
	VoucherID       int
	OperatorID      int
	CompanyID       int
	SubjectID       int
	VoucherRecordID int
}

//CommResp ...
type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}
