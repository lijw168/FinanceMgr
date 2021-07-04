package service

import (
	"context"
	"database/sql"

	"common/log"
)

//RollbackLog ...
func RollbackLog(ctx context.Context, l *log.Logger, funcName string, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		l.ErrorContext(ctx, "[%s] [DB.Rollback:%v]", funcName, err)
	}
}
