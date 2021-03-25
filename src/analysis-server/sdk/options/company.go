package options

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

// type ListCompanyOptions struct {
// 	Filter     []FilterItem `json:"filter"`
// 	Order      []OrderItem  `json:"orders"`
// 	DescOffset int          `json:"desc_offset"`
// 	DescLimit  int          `json:"desc_limit"`
// }

// type DeleteIDOptions struct {
// 	ID int `json:"Id"`
// }

type ModifyCompanyOptions struct {
	CompanyID   int
	CompanyName string
	AbbrevName  string
	Corporator  string
	Phone       string
	Summary     string
	Email       string
	CompanyAddr string
	Backup      string
}
