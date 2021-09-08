package service

import (
	"bytes"
	"common/log"
	"context"
	"database/sql"
	"encoding/json"
)

var GIdInfoService = NewIDInfoService()

//RollbackLog ...
func RollbackLog(ctx context.Context, l *log.Logger, funcName string, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		l.ErrorContext(ctx, "[%s] [DB.Rollback:%v]", funcName, err)
	}
}

func FormatData(srcData interface{}, desData interface{}) error {
	b, err := json.Marshal(srcData)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(desData)
	if err != nil {
		return err
	}
	return nil
}
