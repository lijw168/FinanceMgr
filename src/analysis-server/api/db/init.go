package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// const (
// 	VniPoolBulkSize = 1000
// 	VniBulkSize     = 100
// 	VniUnused       = 0
// 	VniUsed         = 1
// 	VniRelease      = 2
// )

var ErrIpUnavailable = errors.New("Ip address has been used")
var ErrVniUnavailable = errors.New("Vni has been used")

type DbOperator interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DbScanner interface {
	Scan(dest ...interface{}) error
}

func transferSQL(sql string, params []interface{}) (placeholder string) {
	if len(params) > 0 {
		placeholder = "?" + strings.Repeat(",?", len(params)-1)
		return fmt.Sprintf(sql, placeholder)
	}
	return
}

func clause(fields map[string]string, clause map[string]interface{}, sep string) (*string, *[]interface{}, error) {
	buffer := make([]string, 0, 10)
	values := make([]interface{}, 0, 10)
	for k, v := range clause {
		value, ok := fields[k]
		if !ok {
			return nil, nil, errors.New("fields is illegal")
		}
		buffer = append(buffer, fmt.Sprintf(" %v = ? ", value))
		values = append(values, v)
	}

	sql := strings.Join(buffer, sep)
	return &sql, &values, nil
}

//filterNo,the var type is float64, int and string
func transferCountSqlWithNo(table string, filter map[string]interface{}, filterNo map[string]interface{}) (string, []interface{}) {
	strSql := "select count(id) from " + table
	var fk []string
	var fv []interface{}

	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}
	for k, v := range filterNo {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" != ?")
			fv = append(fv, v)
		case []int:
			tmpK += " NOT IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	return strSql, fv
}
func transferCountSql(table string, filter map[string]interface{}) (string, []interface{}) {
	strSql := "select count(id) from " + table
	var fk []string
	var fv []interface{}

	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}

	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	return strSql, fv
}

func transferCountSqlWithFuzzyMatch(table string, filter map[string]interface{}, fuzzyMatchFilter map[string]string) (string, []interface{}) {
	strSql := "select count(id) from " + table
	var fk []string
	var fv []interface{}

	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}
	for k, v := range fuzzyMatchFilter {
		tmpK := camelToUnix(k)
		fk = append(fk, tmpK+" like ?")
		fv = append(fv, v)
	}
	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	return strSql, fv
}

//fuzzyMatchFilter,the value type is  string
func transferListSqlWithFuzzyMatch(table string, filter map[string]interface{}, fuzzyMatchFilter map[string]string,
	field []string, limit int, offset int, order string, od int) (string, []interface{}) {

	fields := strings.Join(field, ",")
	strSql := "select " + fields + " from " + table
	var fk []string
	var fv []interface{}
	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}
	for k, v := range fuzzyMatchFilter {
		tmpK := camelToUnix(k)
		fk = append(fk, tmpK+" like ?")
		fv = append(fv, v)
	}

	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	if order != "" {
		order = camelToUnix(order)
		strSql += " order by " + string(order)
		if od == 1 {
			strSql += " desc"
		}
	}
	if limit >= 0 {
		if offset >= 0 {
			strSql += " LIMIT ?, ?"
			fv = append(fv, offset)
			fv = append(fv, limit)
		} else {
			strSql += " LIMIT ?"
			fv = append(fv, limit)
		}
	}
	return strSql, fv
}

//filterNo,the value type is float64, int and string
func transferListSqlWithNo(table string, filter map[string]interface{}, filterNo map[string]interface{},
	field []string, limit int, offset int, order string, od int) (string, []interface{}) {

	fields := strings.Join(field, ",")
	strSql := "select " + fields + " from " + table
	var fk []string
	var fv []interface{}
	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}
	for k, v := range filterNo {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" != ?")
			fv = append(fv, v)
		case []int:
			tmpK += " NOT IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}

	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	if order != "" {
		order = camelToUnix(order)
		strSql += " order by " + string(order)
		if od == 1 {
			strSql += " desc"
		}
	}
	if limit >= 0 {
		if offset >= 0 {
			strSql += " LIMIT ?, ?"
			fv = append(fv, offset)
			fv = append(fv, limit)
		} else {
			strSql += " LIMIT ?"
			fv = append(fv, limit)
		}
	}
	return strSql, fv
}
func transferListSql(table string, filter map[string]interface{}, field []string, limit int, offset int, order string, od int) (string, []interface{}) {
	fields := strings.Join(field, ",")
	strSql := "select " + fields + " from " + table
	var fk []string
	var fv []interface{}
	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}

	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case string:
			tempS := strings.TrimSpace(v.(string))
			if strings.HasPrefix(tempS, ">") || strings.HasPrefix(tempS, "<") {
				operator := string(tempS[0])
				val := tempS[1:]
				fk = append(fk, fmt.Sprintf("%s %s ?", tmpK, operator))
				fv = append(fv, val)
			} else {
				fk = append(fk, tmpK+" = ?")
				fv = append(fv, v)
			}
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	if order != "" {
		order = camelToUnix(order)
		strSql += " order by " + string(order)
		if od == 1 {
			strSql += " desc"
		}
	}
	if limit >= 0 {
		if offset >= 0 {
			strSql += " LIMIT ?, ?"
			fv = append(fv, offset)
			fv = append(fv, limit)
		} else {
			strSql += " LIMIT ?"
			fv = append(fv, limit)
		}
	}
	return strSql, fv
}

func camelToUnix(s string) string {
	var tmp string
	for i, c := range s {
		if c >= 65 && c <= 90 {
			if i != 0 {
				tmp += "_" + string(c+32)
			} else {
				tmp += string(c + 32)
			}

		} else {
			tmp += string(c)
		}
	}
	return tmp
}

func UpdateDb(do DbOperator, tn string, id string, kv map[string]interface{}) error {
	if len(kv) == 0 {
		return nil
	}
	strSql := "update `" + tn + "`"
	var fk []string
	var fv []interface{}
	for k, v := range kv {
		tmpK := camelToUnix(k)
		switch tmpK {
		case "updated_at":
			fk = append(fk, "`updated_at`=now()")
		case "version":
			fk = append(fk, "`version`=`version`+1")
		default:
			switch v.(type) {
			case float64, int, string:
				fk = append(fk, "`"+tmpK+"`=?")
				fv = append(fv, v)
			default:
				return errors.New("Unsupport value type.")
			}
		}

	}
	strSql += " set " + strings.Join(fk, ",")
	strSql += " where `id`=?"
	fv = append(fv, id)

	_, err := do.Exec(strSql, fv...)
	return err
}

// func listBasicDb(do DbOperator, tn string, filter map[string]interface{}, limit int, offset int, order string, od int) ([]*model.ResourceBasic, error) {
// 	var ret []*model.ResourceBasic
// 	strSql, values := transferListSql(tn, filter, []string{"id", "version"}, limit, offset, order, od)
// 	result, err := do.Query(strSql, values...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer result.Close()
// 	for result.Next() {
// 		r := new(model.ResourceBasic)
// 		err = result.Scan(&(r.Id), &(r.Version))
// 		if err != nil {
// 			return ret, err
// 		}
// 		ret = append(ret, r)
// 	}
// 	return ret, nil
// }

//gen sql with order by list
func transferListSqlWithOrders(table string, filter map[string]interface{}, field []string, limit int, offset int, orders []string, ods []int) (string, []interface{}) {
	fields := strings.Join(field, ",")
	strSql := "select " + fields + " from " + table
	var fk []string
	var fv []interface{}
	handleArrFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else {
				*s += ", ?"
			}
			fv = append(fv, ki)
		}
		return
	}

	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case []int:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]int); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			if vl, ok := v.([]string); ok {
				for _, ki := range vl {
					arr = append(arr, ki)
				}
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v.([]interface{}), &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	if len(filter) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	strSql += genOrderByStr(orders, ods)
	if limit >= 0 {
		if offset >= 0 {
			strSql += " LIMIT ?, ?"
			fv = append(fv, offset)
			fv = append(fv, limit)
		} else {
			strSql += " LIMIT ?"
			fv = append(fv, limit)
		}
	}
	return strSql, fv
}

func genOrderByStr(orders []string, ods []int) string {
	ordStr := ""
	if len(orders) > 0 && len(orders) == len(ods) {
		ordStr += " order by "
		getRangeStr := func(od int) string {
			if od == 0 {
				return "desc"
			} else if od == 1 {
				return "asc"
			} else {
				return ""
			}
		}
		for index, order := range orders {
			ordStr += fmt.Sprintf("%s %s,", order, getRangeStr(ods[index]))
		}
		ordStr = ordStr[:len(ordStr)-1]
	}
	return ordStr
}
