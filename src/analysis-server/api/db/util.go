package db

import (
	"fmt"
)

func GenTableName(iYear int, baseTableName string) string {
	return fmt.Sprintf("%s_%d", baseTableName, iYear)
}
