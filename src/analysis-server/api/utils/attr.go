package utils

import (
	"financeMgr/src/analysis-server/model"
)

type Attribute struct {
	Type int
	Val  Validater
}

func validate_string(d interface{}) bool {
	if _, ok := d.(string); !ok {
		return false
	}
	return true
}

func validate_bool(d interface{}) bool {
	if _, ok := d.(bool); !ok {
		return false
	}
	return true
}

// 由于把json解析成interface{}时，把number解析成float64
// 所以为了判断是否是传过来的整形类型数据，就增加了是否float类型的判断。
// 该函数有bug，需要修改
func validate_int(d interface{}) bool {
	if _, ok := d.(int); !ok {
		validate_float64(d)
	}
	return true
}

func validate_float64(d interface{}) bool {
	if _, ok := d.(float64); !ok {
		return false
	}
	return true
}

func validate_str_arr(d interface{}) bool {
	if l, ok := d.([]interface{}); !ok {
		return false
	} else {
		for _, i := range l {
			if _, ok := i.(string); !ok {
				return false
			}
		}
	}
	return true
}

// 由于把json解析成interface{}时，把number解析成float64
// 所以为了判断是否是传过来的整形类型数据，就增加了是否float类型的判断。该函数有bug，需要修改
func validate_int_arr(d interface{}) bool {
	if l, ok := d.([]interface{}); !ok {
		return false
	} else {
		for _, i := range l {
			if _, ok := i.(int); !ok {
				return false
			}
		}
	}
	return true
}

type Validater func(d interface{}) bool

func ValiFilter(attrs map[string]Attribute, filter []*model.FilterItem) bool {
	for _, f := range filter {
		if f.Field == nil || f.Value == nil {
			return false
		}
		if attr, ok := attrs[*f.Field]; ok {
			if attr.Val != nil {
				if !attr.Val(f.Value) {
					return false
				}
			}
			switch attr.Type {
			case T_Int:
				if !validate_int(f.Value) {
					return false
				}
			case T_Float64:
				if !validate_float64(f.Value) {
					return false
				}
			case T_String:
				if !validate_string(f.Value) {
					return false
				}
			case T_String_Arr:
				if !validate_str_arr(f.Value) {
					return false
				}
			case T_Int_Arr:
				if !validate_int_arr(f.Value) {
					return false
				}
			case T_Bool:
				if !validate_bool(f.Value) {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}
