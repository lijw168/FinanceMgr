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
