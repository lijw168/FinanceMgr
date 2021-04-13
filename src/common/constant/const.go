package constant

import (
	"errors"
	// "time"
)

const (
	Order_Desc = 0
	Order_Asc  = 1
)

var (
	ERR_INVALIDARGUMENT = errors.New("Invalid Argument")
	ERR_NOTEXIST        = errors.New("Not exist")
)
