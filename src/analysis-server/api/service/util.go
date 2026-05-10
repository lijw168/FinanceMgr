package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	dbp "financeMgr/src/analysis-server/api/db"
	"fmt"
)

// RollbackLog ...
func RollbackLog(ctx context.Context, funcName string, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		gLogger.ErrorContext(ctx, "[%s] [DB.Rollback:%v]", funcName, err)
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
func CreateYearVoucherTable(iYear int) error {
	baseTableName := []string{"voucherInfo", "voucherRecordInfo"}
	var err error
	for _, tn := range baseTableName {
		err = createNewTable(tn, dbp.GenTableName(iYear, tn))
		if err != nil {
			gLogger.Error("[CreateYearVoucherTable] [createNewTable: %s]", err.Error())
			return err
		}
	}
	return nil
}

func createNewTable(oldTableName, newTableName string) error {
	gLogger.Debug("[Service/createNewTable] [oldTableName: %s, newTableName: %s]", oldTableName, newTableName)
	//judge ,is not exist
	var c int64
	strSql := "select count(1) from information_schema.TABLES where table_name = ?"
	err := gDb.QueryRow(strSql, newTableName).Scan(&c)
	if err != nil {
		gLogger.Error("[init/service/CreateNewTable] [db.QueryRowContext: %s]", err.Error())
		return err
	}
	if c == 0 {
		//create new table
		strCreateTableSql := fmt.Sprintf("create table %s like %s", newTableName, oldTableName)
		if _, err = gDb.Exec(strCreateTableSql); err != nil {
			gLogger.Error("[init/service/CreateNewTable] [db.Exec: %s]", err.Error())
		}
	}
	return err
}
