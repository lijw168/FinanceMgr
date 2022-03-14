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

//AccSubject ...
type AccSubject struct {
	SubjectID        int
	CompanyID        int
	CommonID         string
	SubjectName      string
	SubjectLevel     int
	SubjectDirection int
	SubjectType      int
	MnemonicCode     string
	SubjectStyle     string
}

type YearBalance struct {
	SubjectID int
	Year      int
	Balance   float64
}

//VoucherInfo ...
type VoucherInfo struct {
	VoucherID int
	CompanyID int
	//oucherYear    int
	VoucherMonth   int
	NumOfMonth     int
	VoucherDate    time.Time
	VoucherFiller  string
	VoucherAuditor string
	BillCount      int
	Status         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//VoucherTemplate ...
type VoucherTemplate struct {
	VoucherTemplateID int
	RefVoucherID      int
	VoucherYear       int
	Illustration      string
	CreatedAt         time.Time
}

//CompanyInfo ...
type CompanyInfo struct {
	CompanyID          int
	CompanyName        string
	AbbrevName         string
	Corporator         string
	Phone              string
	Email              string
	CompanyAddr        string
	Backup             string
	StartAccountPeriod int
	LatestAccountYear  int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CompanyGroupID     int
}

//CompanyGroup ...
type CompanyGroup struct {
	CompanyGroupID int
	GroupName      string
	GroupStatus    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

//IDInfo ...
type IDInfo struct {
	VoucherID         int
	OperatorID        int
	CompanyID         int
	SubjectID         int
	VoucherRecordID   int
	ComGroupID        int
	VoucherTemplateID int
}

//MenuInfo ...
type MenuInfo struct {
	MenuID        int
	MenuName      string
	MenuLevel     int
	ParentMenuID  int
	MenuSerialNum int
}

//CommResp ...
type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}

//AccountOfPeriod
type AccountOfPeriod struct {
	SubjectID         int
	PeriodDebitMoney  float64
	PeriodCreditMoney float64
}

type CalAccuMoney struct {
	CompanyID    int
	SubjectID    int
	VoucherMonth int
	Status       int
	VoucherYear  int
}
