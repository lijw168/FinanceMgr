package options

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
