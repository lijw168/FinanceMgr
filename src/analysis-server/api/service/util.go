package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	dbp "financeMgr/src/analysis-server/api/db"
	"financeMgr/src/common/log"
	"fmt"
)

// RollbackLog ...
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

// CreateYearVoucherTable 创建每年的凭证表，因为进行了分表，所以每年都要创建表
func CreateYearVoucherTable(ctx context.Context, logger *log.Logger, iYear int, db *sql.DB) error {
	baseTableName := []string{"voucherInfo", "voucherRecordInfo"}
	var err error
	for _, tn := range baseTableName {
		err = createNewTable(ctx, logger, db, tn, dbp.GenTableName(iYear, tn))
		if err != nil {
			logger.ErrorContext(ctx, "[CreateYearVoucherTable] [createNewTable: %s]", err.Error())
			return err
		}
	}
	return nil
}

func createNewTable(ctx context.Context, logger *log.Logger, db *sql.DB, oldTableName, newTableName string) error {
	logger.DebugContext(ctx, "[CompanyInfo/db/createNewTable] [oldTableName: %s, newTableName: %s]", oldTableName, newTableName)
	//judge ,is not exist
	var c int64
	strSql := "select count(1) from information_schema.TABLES where table_name = ?"
	err := db.QueryRowContext(ctx, strSql, newTableName).Scan(&c)
	if err != nil {
		logger.ErrorContext(ctx, "[init/service/CreateNewTable] [db.QueryRowContext: %s]", err.Error())
		return err
	}
	if c == 0 {
		//create new table
		strCreateTableSql := fmt.Sprintf("create table %s like %s", newTableName, oldTableName)
		if _, err = db.ExecContext(ctx, strCreateTableSql); err != nil {
			logger.ErrorContext(ctx, "[init/service/CreateNewTable] [db.Exec: %s]", err.Error())
		}
	}
	return err
}
