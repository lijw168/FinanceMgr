package model

type FilterItem struct {
	Field *string     `json:"field"`
	Value interface{} `json:"value"`
}

type OrderItem struct {
	Field     *string `json:"field"`
	Direction *int    `json:"direction"`
}

type BaseParams struct {
	ID *int `json:"id"`
}

// CreateSubjectParams ...
type CreateSubjectParams struct {
	CompanyID        *int    `json:"companyId"`
	CommonID         *string `json:"commonId"`
	SubjectName      *string `json:"subjectName"`
	SubjectLevel     *int    `json:"subjectLevel"`
	SubjectDirection *int    `json:"subjectDirection"`
	SubjectType      *int    `json:"subjectType"`
	MnemonicCode     *string `json:"mnemonicCode"`
	SubjectStyle     *string `json:"subjectStyle"`
}

type ModifySubjectParams struct {
	SubjectID        *int    `json:"subjectId"`
	CompanyID        *int    `json:"companyId"`
	CommonID         *string `json:"commonId"`
	SubjectName      *string `json:"subjectName"`
	SubjectLevel     *int    `json:"subjectLevel"`
	SubjectDirection *int    `json:"subjectDirection"`
	SubjectType      *int    `json:"subjectType"`
	MnemonicCode     *string `json:"mnemonicCode"`
	SubjectStyle     *string `json:"subjectStyle"`
}

type OptYearBalanceParams struct {
	CompanyID *int     `json:"companyId"`
	SubjectID *int     `json:"subjectId"`
	Year      *int     `json:"year"`
	Balance   *float64 `json:"balance"`
	Status    *int     `json:"status"`
}

type OptSubAndBalParam struct {
	SubjectID *int     `json:"subjectId"`
	Balance   *float64 `json:"balance"`
}

type BatchCreateYearBalsParams struct {
	CompanyID     *int                 `json:"companyId"`
	Year          *int                 `json:"year"`
	OptSubAndBals []*OptSubAndBalParam `json:"optSubAndBals"`
}

type BatchDelYearBalsParams struct {
	CompanyID  *int  `json:"companyId"`
	Year       *int  `json:"year"`
	SubjectIDs []int `json:"subjectIds"`
}

type BasicYearBalanceParams struct {
	CompanyID *int `json:"companyId"`
	Year      *int `json:"year"`
	SubjectID *int `json:"subjectId"`
}

// type UpdateItem struct {
// 	Field *string     `json:"field"`
// 	Value interface{} `json:"value"`
// }

type BatchUpdateBalsParams struct {
	CompanyID     *int                 `json:"companyId"`
	Year          *int                 `json:"year"`
	OptSubAndBals []*OptSubAndBalParam `json:"optSubAndBals"`
}

type DeleteSubjectParams struct {
	ID *int `json:"id"`
}

type DescribeIdParams struct {
	ID *int `json:"id"`
}

type ListSubjectParams struct {
	Filter     []*FilterItem `json:"filter"`
	Order      []*OrderItem  `json:"orders"`
	DescOffset *int          `json:"desc_offset"`
	DescLimit  *int          `json:"desc_limit"`
}

type CreateCompanyParams struct {
	CompanyName        *string `json:"companyName"`
	AbbrevName         *string `json:"abbreviationName"`
	Corporator         *string `json:"corporator"`
	Phone              *string `json:"phone"`
	Email              *string `json:"e_mail"`
	CompanyAddr        *string `json:"companyAddr"`
	Backup             *string `json:"backup"`
	StartAccountPeriod *int    `json:"startAccountPeriod"`
}

type ModifyCompanyParams struct {
	CompanyID         *int    `json:"companyId"`
	CompanyName       *string `json:"companyName"`
	AbbrevName        *string `json:"abbreviationName"`
	Corporator        *string `json:"corporator"`
	Phone             *string `json:"phone"`
	Email             *string `json:"e_mail"`
	CompanyAddr       *string `json:"companyAddr"`
	Backup            *string `json:"backup"`
	LatestAccountYear *int    `json:"latestAccountYear"`
}

type ListCompanyParams struct {
	Filter     []*FilterItem `json:"filter"`
	Order      []*OrderItem  `json:"orders"`
	DescOffset *int          `json:"desc_offset"`
	DescLimit  *int          `json:"desc_limit"`
}

type DeleteIDParams struct {
	ID *int `json:"Id"`
}

type IDsParams struct {
	IDs []int `json:"Ids"`
}

type AssociatedCompanyGroupParams struct {
	CompanyGroupID *int  `json:"companyGroupId"`
	CompanyID      *int  `json:"companyId"`
	IsAttach       *bool `json:"isAttach"`
}

type CreateCompanyGroupParams struct {
	GroupName   *string `json:"groupName"`
	GroupStatus *int    `json:"groupStatus"`
}

type ModifyCompanyGroupParams struct {
	CompanyGroupID *int    `json:"companyGroupId"`
	GroupName      *string `json:"groupName"`
	GroupStatus    *int    `json:"groupStatus"`
}

type AuthenInfoParams struct {
	Name      *string `json:"name"`
	Password  *string `json:"password"`
	CompanyID *int    `json:"companyId"`
}

type LoginInfoParams struct {
	Name       *string `json:"name"`
	OperatorID *int    `json:"operatorId"`
	ClientIp   *string `json:"clientIp"`
}

type CreateOptInfoParams struct {
	CompanyID  *int    `json:"companyId"`
	Name       *string `json:"name"`
	Password   *string `json:"password"`
	Job        *string `json:"job"`
	Department *string `json:"department"`
	Status     *int    `json:"Status"`
	Role       *int    `json:"role"`
}

type ModifyOptInfoParams struct {
	OperatorID *int    `json:"operatorId"`
	Name       *string `json:"name"`
	Password   *string `json:"password"`
	Job        *string `json:"job"`
	Department *string `json:"department"`
	Status     *int    `json:"Status"`
	Role       *int    `json:"role"`
}

type ListOperatorsParams struct {
	Filter     []*FilterItem `json:"filter"`
	Order      []*OrderItem  `json:"orders"`
	DescOffset *int          `json:"desc_offset"`
	DescLimit  *int          `json:"desc_limit"`
}

type DescribeNameParams struct {
	Name *string `json:"name"`
}

// type DescribeOperatorParams struct {
// 	Name *string `json:"name"`
// }

type DeleteOperatorParams struct {
	Name *string `json:"name"`
}

type DescribeYearAndIDParams struct {
	VoucherYear *int `json:"voucherYear"`
	ID          *int `json:"id"`
}

type DeleteYearAndIDParams struct {
	VoucherYear *int `json:"voucherYear"`
	ID          *int `json:"id"`
}

type VoucherTemplateParams struct {
	CompanyID    *int    `json:"companyId"`
	RefVoucherID *int    `json:"refVoucherId"`
	VoucherYear  *int    `json:"voucherYear"`
	Illustration *string `json:"illustration"`
}

// VoucherInfoParams ...  BillCount:该参数暂未使用，如果使用时，前端会传入该值。
type VoucherInfoParams struct {
	CompanyID     *int    `json:"companyId"`
	VoucherDate   *int    `json:"voucherDate"`
	VoucherFiller *string `json:"voucherFiller"`
	BillCount     *int    `json:"billCount"`
}

// QueryMonthlyVoucherInfoAttrParameters 查询月的凭证信息的属性参数...
type QueryMonthlyVoucherInfoAttrParameters struct {
	CompanyID    *int `json:"companyId"`
	VoucherYear  *int `json:"voucherYear"`
	VoucherMonth *int `json:"voucherMonth"`
}

// QueryVoucherInfoStatusParams 查询凭证状态...
type QueryVoucherInfoStatusParams struct {
	CompanyID   *int `json:"companyId"`
	VoucherYear *int `json:"voucherYear"`
	Status      *int `json:"status"`
}

type BatchAuditParams struct {
	VoucherYear    *int    `json:"voucherYear"`
	IDs            []int   `json:"ids"`
	Status         *int    `json:"status"`
	VoucherAuditor *string `json:"voucherAuditor"`
}

// ModifyVoucherInfoParams ...
type ModifyVoucherInfoParams struct {
	VoucherYear    *int    `json:"voucherYear"`
	VoucherID      *int    `json:"voucherId"`
	VoucherMonth   *int    `json:"voucherMonth"`
	VoucherDate    *int    `json:"voucherDate"`
	VoucherFiller  *string `json:"voucherFiller"`
	VoucherAuditor *string `json:"voucherAuditor"`
	BillCount      *int    `json:"billCount"`
	Status         *int    `json:"status"`
}

// CreateVoucherRecordParams ...
type CreateVoucherRecordParams struct {
	VoucherID   *int     `json:"voucherId"`
	SubjectName *string  `json:"subjectName"`
	DebitMoney  *float64 `json:"debitMoney"`
	CreditMoney *float64 `json:"creditMoney"`
	Summary     *string  `json:"summary"`
	SubID1      *int     `json:"subId1"`
	// SubID2      *int     `json:"subId2"`
	// SubID3      *int     `json:"subId3"`
	// SubID4      *int     `json:"subId4"`
}

// ModifyVoucherRecordParams ...
type ModifyVoucherRecordParams struct {
	VouRecordID *int     `json:"vouRecordId"`
	SubjectName *string  `json:"subjectName"`
	DebitMoney  *float64 `json:"debitMoney"`
	CreditMoney *float64 `json:"creditMoney"`
	Summary     *string  `json:"summary"`
	SubID1      *int     `json:"subId1"`
	// SubID2      *int     `json:"subId2"`
	// SubID3      *int     `json:"subId3"`
	// SubID4      *int     `json:"subId4"`
}

// CreateVoucherParams...
type CreateVoucherParams struct {
	InfoParams    *VoucherInfoParams           `json:"infoParams"`
	RecordsParams []*CreateVoucherRecordParams `json:"recordsParams"`
}

type UpdateVoucherParams struct {
	VoucherYear         *int                         `json:"voucherYear"`
	ModifyInfoParams    *ModifyVoucherInfoParams     `json:"modifyInfoParams"`
	ModifyRecordsParams []*ModifyVoucherRecordParams `json:"modifyRecordsParams"`
	DelRecordsParams    []int                        `json:"delRecordsParams"`
	AddRecordsParams    []*CreateVoucherRecordParams `json:"addRecordsParams"`
}

//VoucherAuditParams 凭证的审核参数
// type VoucherAuditParams struct {
// 	VoucherID      *int    `json:"voucherId"`
// 	VoucherAuditor *string `json:"voucherAuditor"`
// 	Status         *int    `json:"status"`
// }

// VoucherArrangeParams 凭证整理参数
type VoucherArrangeParams struct {
	VoucherYear       *int  `json:"voucherYear"`
	CompanyID         *int  `json:"companyId"`
	VoucherMonth      *int  `json:"voucherMonth"`
	ArrangeVoucherNum *bool `json:"arrangeVoucherNum"`
}

type ListParams struct {
	Filter     []*FilterItem `json:"filter"`
	Order      []*OrderItem  `json:"orders"`
	DescOffset *int          `json:"desc_offset"`
	DescLimit  *int          `json:"desc_limit"`
}

type ListVoucherInfoParams struct {
	BasicFilter []*FilterItem `json:"basic_filter"`
	AuxiFilter  []*FilterItem `json:"auxi_filter"`
	Order       []*OrderItem  `json:"orders"`
	DescOffset  *int          `json:"desc_offset"`
	DescLimit   *int          `json:"desc_limit"`
}

type DeleteParams struct {
	Name *string `json:"name"`
}

type IDInfoParams struct {
	VoucherID         *int `json:"voucherId"`
	OperatorID        *int `json:"operatorId"`
	CompanyID         *int `json:"companyId"`
	SubjectID         *int `json:"subjectId"`
	VoucherRecordID   *int `json:"voucherRecordId"`
	ComGroupID        *int `json:"comGroupId"`
	VoucherTemplateID *int `json:"voucherTemplateId"`
}

// 以下参数用于报表统计时的参数类型
type CalAccuMoneyParams struct {
	CompanyID   *int `json:"companyId"`
	SubjectID   *int `json:"subjectId"`
	VoucherDate *int `json:"VoucherDate"`
	Status      *int `json:"status"`
	VoucherYear *int `json:"voucherYear"`
}

type BatchCalAccuMoneyParams struct {
	CompanyID    *int  `json:"companyId"`
	SubjectIDArr []int `json:"subId1Arr"`
	VoucherDate  *int  `json:"VoucherDate"`
	Status       *int  `json:"status"`
	VoucherYear  *int  `json:"voucherYear"`
}

type CalAmountOfPeriodParams struct {
	CompanyID    *int  `json:"companyId"`
	SubjectIDArr []int `json:"subId1Arr"`
	StartMonth   *int  `json:"startMonth"`
	EndMonth     *int  `json:"endMonth"`
	Status       *int  `json:"status"`
	VoucherYear  *int  `json:"voucherYear"`
}
