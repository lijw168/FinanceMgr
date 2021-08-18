package constant

import (
	"errors"
	// "time"
)

const (
	Order_Asc  = 0
	Order_Desc = 1
)

var (
	ERR_INVALIDARGUMENT = errors.New("Invalid Argument")
	ERR_NOTEXIST        = errors.New("Not exist")
)
