package options

//VoucherInfoOptions ...
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

// type ListOptions struct {
// 	Filter     []FilterItem
// 	Order      []OrderItem
// 	DescOffset int
// 	DescLimit  int
// }

// type DeleteOptions struct {
// 	Name string
// }

// type IDInfoOptions struct {
// 	VoucherID       int
// 	CompanyID       int
// 	SubjectID       int
// 	VoucherRecordID int
// }
