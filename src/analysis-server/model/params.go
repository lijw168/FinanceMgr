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
	SubjectName  *string `json:"subjectName"`
	SubjectLevel *int    `json:"subjectLevel"`
}

type ModifySubjectParams struct {
	SubjectID    *int    `json:"subjectId"`
	SubjectName  *string `json:"subjectName"`
	SubjectLevel *int    `json:"subjectLevel"`
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
	CompanyName *string `json:"companyName"`
	AbbrevName  *string `json:"abbreviationName"`
	Corporator  *string `json:"corporator"`
	Phone       *string `json:"phone"`
	Email       *string `json:"e_mail"`
	CompanyAddr *string `json:"companyAddr"`
	Backup      *string `json:"backup"`
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

type ModifyCompanyParams struct {
	CompanyID   *int    `json:"companyId"`
	CompanyName *string `json:"companyName"`
	AbbrevName  *string `json:"abbreviationName"`
	Corporator  *string `json:"corporator"`
	Phone       *string `json:"phone"`
	Email       *string `json:"e_mail"`
	CompanyAddr *string `json:"companyAddr"`
	Backup      *string `json:"backup"`
}

type AuthenInfoParams struct {
	Name      *string `json:"name"`
	Password  *string `json:"password"`
	CompanyID *int    `json:"companyId"`
}

type LoginInfoParams struct {
	Name *string `json:"name"`
	//Status   *int    `json:"status"`
	ClientIp *string `json:"clientIp"`
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

//VoucherInfoParams ...
type VoucherInfoParams struct {
	CompanyID    *int `json:"companyId"`
	VoucherMonth *int `json:"voucherMonth"`
}

//CreateVoucherRecordParams ...
type CreateVoucherRecordParams struct {
	VoucherID   *int     `json:"voucherId"`
	SubjectName *string  `json:"subjectName"`
	DebitMoney  *float64 `json:"debitMoney"`
	CreditMoney *float64 `json:"creditMoney"`
	Summary     *string  `json:"summary"`
	SubID1      *int     `json:"subId1"`
	SubID2      *int     `json:"subId2"`
	SubID3      *int     `json:"subId3"`
	SubID4      *int     `json:"subId4"`
	BillCount   *int     `json:"billCount"`
}

//ModifyVoucherRecordParams ...
type ModifyVoucherRecordParams struct {
	VouRecordID *int     `json:"vouRecordId"`
	SubjectName *string  `json:"subjectName"`
	DebitMoney  *float64 `json:"debitMoney"`
	CreditMoney *float64 `json:"creditMoney"`
	Summary     *string  `json:"summary"`
	SubID1      *int     `json:"subId1"`
	SubID2      *int     `json:"subId2"`
	SubID3      *int     `json:"subId3"`
	SubID4      *int     `json:"subId4"`
	BillCount   *int     `json:"billCount"`
}

//VoucherParams...
type VoucherParams struct {
	InfoParams    *VoucherInfoParams           `json:"infoParams"`
	RecordsParams []*CreateVoucherRecordParams `json:"recordsParams"`
}

type ListParams struct {
	Filter     []*FilterItem `json:"filter"`
	Order      []*OrderItem  `json:"orders"`
	DescOffset *int          `json:"desc_offset"`
	DescLimit  *int          `json:"desc_limit"`
}

type DeleteParams struct {
	Name *string `json:"name"`
}

type IDInfoParams struct {
	VoucherID       *int `json:"voucherID"`
	CompanyID       *int `json:"companyID"`
	SubjectID       *int `json:"subjectID"`
	VoucherRecordID *int `json:"voucher_recordID"`
}
