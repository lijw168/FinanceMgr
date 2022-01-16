package options

//common information
type CommResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseOptions struct {
	ID int
}

type ListOptions struct {
	Filter map[string]interface{}
	Offset int
	Limit  int
	Orders map[string]int
}

//login information
type AuthenInfoOptions struct {
	Name      string
	Password  string
	CompanyID int
}

type LogoutOptions struct {
	Name        string
	AccessToken string
}

//operator information
type CreateOptInfoOptions struct {
	CompanyID  int
	Name       string
	Password   string
	Job        string
	Department string
	Status     int
	Role       int
}

type ModifyOptInfoOptions struct {
	OperatorID int
	Name       string
	Password   string
	Job        string
	Department string
	Status     int
	Role       int
}

type NameOptions struct {
	Name string
}

type QueryMaxNumOfMonthOption struct {
	CompanyID    int
	VoucherYear  int
	VoucherMonth int
}

//voucher begin;VoucherInfoOptions ...
type VoucherInfoOptions struct {
	CompanyID     int
	VoucherMonth  int
	VoucherFiller string
}

//ModifyVoucherInfoOptions ...
type ModifyVoucherInfoOptions struct {
	VoucherYear    int
	VoucherID      int
	VoucherMonth   int
	VoucherDate    int
	VoucherFiller  string
	VoucherAuditor string
	BillCount      int
	Status         int
}

//VoucherRecordOptions ...
type CreateVoucherRecordOptions struct {
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
}

//ModifyVoucherRecordOptions ...
type ModifyVoucherRecordOptions struct {
	VouRecordID int
	SubjectName string
	DebitMoney  float64
	CreditMoney float64
	Summary     string
	SubID1      int
	SubID2      int
	SubID3      int
	SubID4      int
	BillCount   int
	Status      int
}

//VoucherOptions...
type VoucherOptions struct {
	InfoOptions    VoucherInfoOptions
	RecordsOptions []CreateVoucherRecordOptions
}

//VoucherArrangeOptions 凭证整理参数
type VoucherArrangeOptions struct {
	CompanyID         int
	VoucherYear       int
	VoucherMonth      int
	ArrangeVoucherNum bool
}

type DescribeYearAndIDOptions struct {
	VoucherYear int
	ID          int
}

type DeleteYearAndIDOptions struct {
	VoucherYear int
	ID          int
}

//voucher end

//commpany option
type CreateCompanyOptions struct {
	CompanyName        string
	AbbrevName         string
	Corporator         string
	Phone              string
	Summary            string
	Email              string
	CompanyAddr        string
	Backup             string
	StartAccountPeriod int
}

type ModifyCompanyOptions struct {
	CompanyID         int
	CompanyName       string
	AbbrevName        string
	Corporator        string
	Phone             string
	Email             string
	CompanyAddr       string
	Backup            string
	LatestAccountYear int
}

type AssociatedCompanyGroupOptions struct {
	CompanyGroupID int
	CompanyID      int
	IsAttach       bool
}

type CreateCompanyGroupOptions struct {
	GroupName   string
	GroupStatus int
}

type ModifyCompanyGroupOptions struct {
	CompanyGroupID int
	GroupName      string
	GroupStatus    int
}

// account subject option ...
type CreateSubjectOptions struct {
	CompanyID        int
	CommonID         string
	SubjectName      string
	SubjectLevel     int
	SubjectDirection int
	SubjectType      int
	MnemonicCode     string
	SubjectStyle     string
}

type ModifySubjectOptions struct {
	SubjectID        int
	CompanyID        int
	CommonID         string
	SubjectName      string
	SubjectLevel     int
	SubjectDirection int
	SubjectType      int
}

//YearBalanceOption ...
type YearBalanceOption struct {
	SubjectID        int
	Summary          string
	SubjectDirection int
	Balance          float64
}
