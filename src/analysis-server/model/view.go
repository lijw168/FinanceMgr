package model

// api response data model

import (
	"time"
)

//AccSubjectView ...
type AccSubjectView struct {
	SubjectID    string `json:"subjectId"`
	SubjectName  string `json:"subjectName"`
	SubjectLevel int    `json:"subjectLevel"`
}

type CompanyView struct {
	CompanyID   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
	AbbrevName  string `json:"abbreviationName"`
	Corporator  string `json:"corporator"`
	Phone       string `json:"phone"`
	Summary     string `json:"summary"`
	Email       string `json:"e_mail"`
	CompanyAddr string `json:"companyAddr"`
	Backup      string `json:"backup"`
}

type OperatorInfoView struct {
	CompanyID  int    `json:"companyId"`
	Name       string `json:"name"`
	Job        string `json:"job"`
	Department string `json:"department"`
	Status     int    `json:"Status"`
	Role       int    `json:"role"`
}

type VoucherInfoView struct {
	VoucherID    int       `json:"voucherId"`
	CompanyID    int       `json:"companyId"`
	VoucherMonth int       `json:"voucherMonth"`
	NumOfMonth   int       `json:"numOfMonth"`
	VoucherDate  time.Time `json:"voucherDate"`
}

//VoucherRecordView ...
type VoucherRecordView struct {
	RecordID    int     `json:"recordId"`
	VoucherID   int     `json:"voucherId"`
	SubjectName string  `json:"subjectName"`
	DebitMoney  float64 `json:"debitMoney"`
	CreditMoney float64 `json:"creditMoney"`
	Summary     string  `json:"summary"`
	SubID1      int     `json:"subId1"`
	SubID2      int     `json:"subId2"`
	SubID3      int     `json:"subId3"`
	SubID4      int     `json:"subId4"`
	BillCount   int     `json:"billCount"`
}

//表示符合条件的voucher records的总数，通过该变量，客户端用于判断，是否再次获取voucher records
type VoucherView struct {
	VouInfoView         VoucherInfoView     `json:"voucherInfoView"`
	VouRecordTotalCount int                 `json:"vouRecordCount"`
	VouRecordViewSli    []VoucherRecordView `json:"voucherRecordView"`
}

type IDInfoView struct {
	VoucherID       int `json:"voucherID"`
	CompanyID       int `json:"companyID"`
	SubjectID       int `json:"subjectID"`
	VoucherRecordID int `json:"voucher_recordID"`
}
