package options

/*
type CreateHostOptions struct {
	Name        string
	MgmtIp      string
	DataIp      string
	AdminStatus int
	Rack        string
}
*/

type ModifyAdminStatusOptions struct {
	Name        string
	AdminStatus string
}

type ModifyStatusOptions struct {
	Name   string
	Status int
	Force  bool
}

type DescribeHostOptions struct {
	Id string
}

type DeleteHostOptions struct {
	Id string
}
