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

//voucher begin;VoucherInfoOptions ...
type VoucherInfoOptions struct {
	CompanyID    int
	VoucherMonth int
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
}

//VoucherOptions...
type VoucherOptions struct {
	InfoOptions    VoucherInfoOptions
	RecordsOptions []CreateVoucherRecordOptions
}

//voucher end

//commpany option
type CreateCompanyOptions struct {
	CompanyName string
	AbbrevName  string
	Corporator  string
	Phone       string
	Summary     string
	Email       string
	CompanyAddr string
	Backup      string
}

type ModifyCompanyOptions struct {
	CompanyID   int
	CompanyName string
	AbbrevName  string
	Corporator  string
	Phone       string
	Email       string
	CompanyAddr string
	Backup      string
}

// account subject option ...
type CreateSubjectOptions struct {
	CommonID     string
	SubjectName  string
	SubjectLevel int
}

type ModifySubjectOptions struct {
	SubjectID    int
	CommonID     string
	SubjectName  string
	SubjectLevel int
}
