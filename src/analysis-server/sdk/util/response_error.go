package util

type RespErr struct {
	Code int
	Err  error
}

func (re *RespErr) Error() string {
	if re == nil {
		return "<nil>"
	}
	return re.Err.Error()
}
