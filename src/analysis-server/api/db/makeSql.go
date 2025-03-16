package db

import (
	"context"
	"database/sql"
	"errors"
	"financeMgr/src/analysis-server/model"
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

// filterNo,the var type is float64, int and string
func transferCountSqlWithNo(table string, filter map[string]interface{}, filterNo map[string]interface{}) (string, []interface{}) {
	//strSql := "select count(id) from " + table
	strSql := "select count(*) from " + table
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

// support multi-condition:accurate match,!= ,not,in,like,between ... and,
// fuzzyMatchFilter,the value type is  string
// intervalFilter,the value type is numerical value
// filterNo,the value type is float64, int and string
// filter,support '<' ,'>','=',but,"<,>",the value type is string
// sort multi columns
// explain: usse makeListSqlWithMultiCondition ,replace the function
func transferListSqlWithMutiCondition(table string, filterNo map[string]interface{}, filter map[string]interface{},
	intervalFilter map[string]interface{}, fuzzyMatchFilter map[string]string, field []string,
	limit int, offset int, order string, od int) (string, []interface{}) {
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
	handleBetweenFilter := func(arr []interface{}, s *string) (fv []interface{}) {
		for i, ki := range arr {
			if i == 0 {
				*s += "?"
			} else if i == 1 {
				*s += " and ? "
			} else {
				panic("the between parameter is valid")
			}
			fv = append(fv, ki)
		}
		return
	}
	for k, v := range filterNo {
		tmpK := camelToUnix(k)
		switch v := v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" != ?")
			fv = append(fv, v)
		case []int:
			tmpK += " NOT IN ("
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		}
	}
	for k, v := range filter {
		tmpK := camelToUnix(k)
		switch v := v.(type) {
		case float64, int:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case string:
			tempS := strings.TrimSpace(v)
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
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpFv := handleArrFilter(arr, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []interface{}:
			tmpK += " IN ("
			tmpFv := handleArrFilter(v, &tmpK)
			tmpK += ")"
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		default:
			panic("invalid parameter type")
		}
	}
	for k, v := range intervalFilter {
		tmpK := camelToUnix(k)
		switch v := v.(type) {
		case []int:
			tmpK += " between "
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpFv := handleBetweenFilter(arr, &tmpK)
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		case []float64:
			tmpK += " between "
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpFv := handleBetweenFilter(arr, &tmpK)
			fv = append(fv, tmpFv...)
			fk = append(fk, tmpK)
		default:
			panic("invalid parameter type")
		}
	}
	for k, v := range fuzzyMatchFilter {
		tmpK := camelToUnix(k)
		fk = append(fk, tmpK+" like ?")
		fv = append(fv, v)
	}
	if len(fk) > 0 {
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

// filterNo,the value type is float64, int and string
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

func transferDeleteSql(table string, filter map[string]interface{}) (string, []interface{}) {
	strSql := "delete  from " + table
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

// gen sql with order by list
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

func processCommonConditionFilter(conditionFilter []interface{}, sql *string) (filterValue []interface{}) {
	for i, fv := range conditionFilter {
		if i == 0 {
			*sql += "?"
		} else {
			*sql += ", ?"
		}
		filterValue = append(filterValue, fv)
	}
	return
}

func processBetweenFilter(btConditionFilter []interface{}, sql *string) (filterValue []interface{}) {
	for i, kv := range btConditionFilter {
		if i == 0 {
			*sql += "?"
		} else if i == 1 {
			*sql += " and ? "
		} else {
			panic("the between parameter is valid")
		}
		filterValue = append(filterValue, kv)
	}
	return
}

func addCommonFilterCondition(filter map[string]interface{}) (fk []string, fv []interface{}) {
	var tmpK string
	for k, v := range filter {
		tmpK = camelToUnix(k)
		switch v := v.(type) {
		case float64, int:
			fk = append(fk, tmpK+" = ?")
			fv = append(fv, v)
		case string:
			tempS := strings.TrimSpace(v)
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
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpValue := processCommonConditionFilter(arr, &tmpK)
			tmpK += ")"
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		case []string:
			tmpK += " IN ("
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpValue := processCommonConditionFilter(arr, &tmpK)
			tmpK += ")"
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		case []interface{}:
			tmpK += " IN ("
			tmpValue := processCommonConditionFilter(v, &tmpK)
			tmpK += ")"
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		default:
			panic("invalid parameter type")
		}
	}
	return
}

func addNoFilterCondition(filterNo map[string]interface{}) (fk []string, fv []interface{}) {
	var tmpK string
	for k, v := range filterNo {
		tmpK = camelToUnix(k)
		switch v := v.(type) {
		case float64, int, string:
			fk = append(fk, tmpK+" != ?")
			fv = append(fv, v)
		case []int:
			tmpK += " NOT IN ("
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpValue := processCommonConditionFilter(arr, &tmpK)
			tmpK += ")"
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		case []interface{}:
			tmpK += " NOT IN ("
			tmpValue := processCommonConditionFilter(v, &tmpK)
			tmpK += ")"
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		}
	}
	return
}

func addBetweenFilterSql(intervalFilter map[string]interface{}) (fk []string, fv []interface{}) {
	var tmpK string
	for k, v := range intervalFilter {
		tmpK = camelToUnix(k)
		switch v := v.(type) {
		case []int:
			tmpK += " between "
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpValue := processBetweenFilter(arr, &tmpK)
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		case []float64:
			tmpK += " between "
			arr := []interface{}{}
			for _, ki := range v {
				arr = append(arr, ki)
			}
			tmpValue := processBetweenFilter(arr, &tmpK)
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		case []interface{}:
			tmpK += " between "
			tmpValue := processBetweenFilter(v, &tmpK)
			fk = append(fk, tmpK)
			fv = append(fv, tmpValue...)
		default:
			panic("invalid parameter type")
		}
	}
	return
}

func addFuzzyMatchFilter(fuzzyMatchFilter map[string]string) (fk []string, fv []interface{}) {
	var tmpK string
	for k, v := range fuzzyMatchFilter {
		tmpK = camelToUnix(k)
		fk = append(fk, tmpK+` like "%?%" `)
		fv = append(fv, v)
	}
	return
}

// order by 后面的字段的顺序要和传入的字段的顺序一致，否则会出现错误。所以修改函数参数为slice
func addOrderSql(orderFiler []*model.OrderItem) string {
	if len(orderFiler) == 0 {
		return ""
	}
	var tmpK string
	iCount := 0
	orderSql := " order by "
	for _, v := range orderFiler {
		tmpK = camelToUnix(*v.Field)
		if *v.Direction == 1 {
			if iCount == 0 {
				orderSql += fmt.Sprintf("%s desc ", tmpK)
			} else {
				orderSql += fmt.Sprintf(", %s desc ", tmpK)
			}

		} else {
			if iCount == 0 {
				orderSql += fmt.Sprintf("%s ", tmpK)
			} else {
				orderSql += fmt.Sprintf(", %s ", tmpK)
			}
		}
		iCount++
	}
	return orderSql
}

func addLimitSql(limit, offset int) (strSql string, fv []interface{}) {
	if limit > 0 {
		if offset > 0 {
			strSql = " LIMIT ?, ?"
			fv = append(fv, offset, limit)
		} else {
			strSql = " LIMIT ?"
			fv = append(fv, limit)
		}
	}
	return
}

// support multi-condition:accurate match,!= ,not,in,like,between ... and,
// fuzzyMatchFilter,the value type is  string
// intervalFilter,the value type is numerical value
// filterNo,the value type is float64, int and string
// filter,support '<' ,'>','=',but,"<,>",the value type is string
// sort multi columns
func makeListSqlWithMultiCondition(table string, field []string,
	filterNo map[string]interface{}, filter map[string]interface{},
	intervalFilter map[string]interface{}, fuzzyMatchFilter map[string]string,
	orderFiler []*model.OrderItem, limit int, offset int) (string, []interface{}) {

	fields := strings.Join(field, ",")
	strSql := "select " + fields + " from " + table
	var fk, tmpKSlice []string
	var fv, tmpVSlice []interface{}
	if len(filterNo) > 0 {
		tmpKSlice, tmpVSlice := addNoFilterCondition(filterNo)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(filter) > 0 {
		tmpKSlice, tmpVSlice = addCommonFilterCondition(filter)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(intervalFilter) > 0 {
		tmpKSlice, tmpVSlice = addBetweenFilterSql(intervalFilter)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(fuzzyMatchFilter) > 0 {
		tmpKSlice, tmpVSlice = addFuzzyMatchFilter(fuzzyMatchFilter)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(fk) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	if len(orderFiler) > 0 {
		strSql += addOrderSql(orderFiler)
	}
	var tmpSql string
	tmpSql, tmpVSlice = addLimitSql(limit, offset)
	strSql += tmpSql
	fv = append(fv, tmpVSlice...)
	return strSql, fv
}

func transferUpdateSql(table string, filter map[string]interface{}, updateField map[string]interface{}) (string, []interface{}) {
	strSql := "update " + table + " set "
	var fv []interface{}
	var first bool = true
	for key, value := range updateField {
		dbKey := camelToUnix(key)
		if first {
			strSql += dbKey + "=?"
			first = false
		} else {
			strSql += "," + dbKey + "=?"
		}
		fv = append(fv, value)
	}
	if first {
		return "", nil
	}
	var fk []string

	tmpKSlice, tmpVSlice := addCommonFilterCondition(filter)
	if len(tmpKSlice) > 0 {
		fk = append(fk, tmpKSlice...)
	}
	if len(tmpVSlice) > 0 {
		fv = append(fv, tmpVSlice...)
	}

	if len(tmpKSlice) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	return strSql, fv
}

func transferCountSql(table string, filter map[string]interface{}) (string, []interface{}) {
	//strSql := "select count(id) from " + table
	strSql := "select count(*) from " + table
	var fk []string
	var fv []interface{}

	tmpKSlice, tmpVSlice := addCommonFilterCondition(filter)
	if len(tmpKSlice) > 0 {
		fk = append(fk, tmpKSlice...)
	}
	if len(tmpVSlice) > 0 {
		fv = append(fv, tmpVSlice...)
	}
	if len(fk) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	return strSql, fv
}

// support multi-condition:accurate match,!= ,not,in,like,between ... and,
// fuzzyMatchFilter,the value type is  string
// intervalFilter,the value type is numerical value
// filterNo,the value type is float64, int and string
// filter,support '<' ,'>','=',but,"<,>",the value type is string
func makeUpdateSqlWithMultiCondition(table string, updateField map[string]interface{},
	filterNo map[string]interface{}, filter map[string]interface{},
	intervalFilter map[string]interface{},
	fuzzyMatchFilter map[string]string) (string, []interface{}) {

	strSql := "update " + table + " set "
	var fv, tmpVSlice []interface{}
	var first bool = true
	for key, value := range updateField {
		dbKey := camelToUnix(key)
		if first {
			strSql += dbKey + "=?"
			first = false
		} else {
			strSql += "," + dbKey + "=?"
		}
		fv = append(fv, value)
	}
	if first {
		return "", nil
	}
	var fk, tmpKSlice []string
	if len(filterNo) > 0 {
		tmpKSlice, tmpVSlice := addNoFilterCondition(filterNo)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(filter) > 0 {
		tmpKSlice, tmpVSlice = addCommonFilterCondition(filter)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(intervalFilter) > 0 {
		tmpKSlice, tmpVSlice = addBetweenFilterSql(intervalFilter)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(fuzzyMatchFilter) > 0 {
		tmpKSlice, tmpVSlice = addFuzzyMatchFilter(fuzzyMatchFilter)
		if len(tmpKSlice) > 0 {
			fk = append(fk, tmpKSlice...)
		}
		if len(tmpVSlice) > 0 {
			fv = append(fv, tmpVSlice...)
		}
	}
	if len(fk) > 0 {
		strSql += " where " + strings.Join(fk, " and ")
	}
	return strSql, fv
}
