package model

// api response data model

import (
	"time"
)

// 该结构体里的Tc字段是为了满足分页查询时，客户端需要知道符合条件的总记录数的需求而设计的，Elements字段才是实际的返回数据
// 例如，ListCompany接口需要分页查询公司列表，客户端需要知道符合条件的公司总数，以判断是否需要继续分页查询，所以在ListCompany接口中，返回的数据结构就是DescData类型的，其中Tc字段表示符合条件的公司总数，Elements字段才是实际的公司列表数据。
// 但在2026.4.18之前的设计中，都没有继续分页的判断，以后要进行修改。
type DescData struct {
	Tc       int64       `json:"total_count"`
	Elements interface{} `json:"elements"`
}

// AccSubjectView ...
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

// YearBalance ...
type YearBalanceView struct {
	Year      int     `json:"year"`
	SubjectID int     `json:"subjectId"`
	Balance   float64 `json:"balance"`
	Status    byte    `json:"status"`
}

type CompanyView struct {
	CompanyID   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
	AbbrevName  string `json:"abbreviationName"`
	Corporator  string `json:"corporator"`
	Phone       string `json:"phone"`
	Email       string `json:"e_mail"`
	CompanyAddr string `json:"companyAddr"`
	Backup      string `json:"backup"`
	// StartAccountPeriod这个字段是为了满足客户端的需求而设计的，表示公司的起始建账的时间，包括年和月，不包括具体哪天，
	// 例如201901，202001等，客户端通过该字段可以知道公司的建账期间，以便于进行凭证的查询和统计等操作。
	BeginAccountDate  int       `json:"beginAccountDate"`
	LatestAccountYear int       `json:"latestAccountYear"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CompanyGroupID    int       `json:"companyGroupId"`
}

type CompanyGroupView struct {
	CompanyGroupID int       `json:"companyGroupId"`
	GroupName      string    `json:"groupName"`
	GroupStatus    int       `json:"groupStatus"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// LoginInfoView ...
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

// VoucherRecordView ...
type VoucherRecordView struct {
	RecordID    int     `json:"recordId"`
	VoucherID   int     `json:"voucherId"`
	SubjectName string  `json:"subjectName"`
	DebitMoney  float64 `json:"debitMoney"`
	CreditMoney float64 `json:"creditMoney"`
	Summary     string  `json:"summary"`
	SubID1      int     `json:"subId1"`
	RecordPos   int     `json:"recordPos"`
	// SubID3      int     `json:"subId3"`
	// SubID4      int     `json:"subId4"`
}

// 表示符合条件的voucher records的总数，通过该变量，客户端用于判断，是否再次获取voucher records
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

type CompanyAccountYearInfoView struct {
	CompanyId         int    `json:"companyId"`
	CompanyName       string `json:"companyName"`
	BeginAccountYear  int    `json:"beginAccountYear"`
	LatestAccountYear int    `json:"latestAccountYear"`
}

// MenuInfo ...
type MenuInfoView struct {
	MenuID       int    `json:"menuId"`
	MenuName     string `json:"menuName"`
	MenuLevel    int    `json:"menuLevel"`
	ParentMenuID int    `json:"parentMenuId"`
}

// AccuMoneyValue ...
type AccuMoneyValueView struct {
	SubjectID       int     `json:"subjectId"`
	AccuDebitMoney  float64 `json:"accuDebitMoney"`
	AccuCreditMoney float64 `json:"accuCreditMoney"`
}

//AccountOfPeriod
// type AccountOfPeriodView struct {
// 	SubjectID            int     `json:"subjectId"`
// 	PeriodDebitMoneySum  float64 `json:"periodDebitMoneySum"`
// 	PeriodCreditMoneySum float64 `json:"periodCreditMoneySum"`
// }
