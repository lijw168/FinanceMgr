package model

// api response data model

import (
	"time"
)

type DescData struct {
	Tc       int64       `json:"total_count"`
	Elements interface{} `json:"elements"`
}

//AccSubjectView ...
type AccSubjectView struct {
	SubjectID        int    `json:"subjectId"`
	CompanyID        int    `json:"companyId"`
	CommonID         string `json:"commonId"`
	SubjectName      string `json:"subjectName"`
	SubjectLevel     int    `json:"subjectLevel"`
	SubjectDirection int    `json:"subjectDirection"`
	SubjectType      int    `json:"subjectType"`
	MnemonicCode     string `json:"mnemonicCode"`
	SubjectStyle     string `json:"subjectStyle"`
}

type CompanyView struct {
	CompanyID          int       `json:"companyId"`
	CompanyName        string    `json:"companyName"`
	AbbrevName         string    `json:"abbreviationName"`
	Corporator         string    `json:"corporator"`
	Phone              string    `json:"phone"`
	Email              string    `json:"e_mail"`
	CompanyAddr        string    `json:"companyAddr"`
	Backup             string    `json:"backup"`
	StartAccountPeriod int       `json:"startAccountPeriod"`
	LatestAccountYear  int       `json:"latestAccountPeriod"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	CompanyGroupID     int       `json:"companyGroupId"`
}

type CompanyGroupView struct {
	CompanyGroupID int       `json:"companyGroupId"`
	GroupName      string    `json:"groupName"`
	GroupStatus    int       `json:"groupStatus"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

//LoginInfoView ...
type LoginInfoView struct {
	OperatorID  int       `json:"operatorId"`
	Name        string    `json:"name"`
	Status      int       `json:"status"`
	ClientIp    string    `json:"clientIp"`
	BeginedAt   time.Time `json:"beginedAt"`
	EndedAt     time.Time `json:"endedAt"`
	AccessToken string    `json:"accessToken"`
}

type StatusCheckoutView struct {
	OperatorID int    `json:"operatorId"`
	Name       string `json:"name"`
	Status     int    `json:"status"`
}

type OperatorInfoView struct {
	OperatorID int       `json:"operatorId"`
	CompanyID  int       `json:"companyId"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	Job        string    `json:"job"`
	Department string    `json:"department"`
	Status     int       `json:"status"`
	Role       int       `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type VoucherInfoView struct {
	VoucherID      int    `json:"voucherId"`
	CompanyID      int    `json:"companyId"`
	VoucherMonth   int    `json:"voucherMonth"`
	NumOfMonth     int    `json:"numOfMonth"`
	VoucherDate    int    `json:"voucherDate"`
	VoucherFiller  string `json:"voucherFiller"`
	VoucherAuditor string `json:"voucherAuditor"`
	BillCount      int    `json:"billCount"`
	Status         int    `json:"status"`
}

type VoucherTemplateView struct {
	VoucherTemplateID int       `json:"voucherTemplateId"`
	RefVoucherID      int       `json:"refVoucherId"`
	VoucherYear       int       `json:"voucherYear"`
	Illustration      string    `json:"illustration"`
	CreatedAt         time.Time `json:"createdAt"`
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
}

//表示符合条件的voucher records的总数，通过该变量，客户端用于判断，是否再次获取voucher records
type VoucherView struct {
	VouInfoView         VoucherInfoView     `json:"voucherInfoView"`
	VouRecordTotalCount int                 `json:"vouRecordCount"`
	VouRecordViewSli    []VoucherRecordView `json:"voucherRecordView"`
}

type IDInfoView struct {
	VoucherID         int `json:"voucherId"`
	OperatorID        int `json:"operatorId"`
	CompanyID         int `json:"companyId"`
	SubjectID         int `json:"subjectId"`
	VoucherRecordID   int `json:"voucherRecordId"`
	ComGroupID        int `json:"comGroupId"`
	VoucherTemplateID int `json:"voucherTemplateId"`
}

type ResourceInfoView struct {
	CompanyId   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
	YearSlice   []int  `json:"year"`
}

//MenuInfo ...
type MenuInfoView struct {
	MenuID       int    `json:"menuId"`
	MenuName     string `json:"menuName"`
	MenuLevel    int    `json:"menuLevel"`
	ParentMenuID int    `json:"parentMenuId"`
}
