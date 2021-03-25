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
// CreateSubjectOptions ...
type CreateSubjectOptions struct {
	SubjectName  string
	SubjectLevel int
}

type ModifySubjectOptions struct {
	SubjectID    int
	SubjectName  string
	SubjectLevel int
}

// type DeleteSubjectOptions struct {
// 	ID int
// }

// type DescribeIdOptions struct {
// 	ID int
// }

// type ListSubjectOptions struct {
// 	Filter []*FilterItem `json:"filter"`
// }
