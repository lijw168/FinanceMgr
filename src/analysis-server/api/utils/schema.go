package utils

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"unicode/utf8"
)

const (
	T_String = iota
	T_Int
	T_Bool
	T_Uuid
	T_Struct
	T_Filter
	T_Order
	T_Array
	T_Text
	T_IP
	T_UuidNull
	T_String_Arr
	T_Int_Arr
)

var (
	uuidP     string = `^[a-z0-9]([a-z0-9-])*$`
	textP     string = `^([a-zA-Z0-9\_-]|[\p{Han}])*$`
	textRe    *regexp.Regexp
	uuidRe    *regexp.Regexp
	schemaVal *SchemaValidator = &SchemaValidator{}
)

type ValidatorFunc func(d interface{}) bool

type SchemaItem struct {
	Type     int
	Val      ValidatorFunc
	Required bool
	SubItem  interface{}
}

type SchemaValidator struct{}

func (s *SchemaValidator) validate_string(d interface{}) bool {
	_, ok := d.(string)
	return ok
}

func (s *SchemaValidator) validate_bool(d interface{}) bool {
	_, ok := d.(bool)
	return ok
}

func (s *SchemaValidator) validate_int(d interface{}) bool {
	_, ok := d.(float64)
	return ok
}

func (s *SchemaValidator) validate_uuid(d interface{}) bool {
	if uuidRe == nil {
		uuidRe = regexp.MustCompile(uuidP)
	}
	if id, ok := d.(string); ok {
		return uuidRe.Match([]byte(id))
	}
	return false
}

func (s *SchemaValidator) validate_uuidNull(d interface{}) bool {
	if id, ok := d.(string); ok {
		return id == "" || s.validate_uuid(d)
	}
	return false
}

func (s *SchemaValidator) validate_text(d interface{}) bool {
	if textRe == nil {
		textRe = regexp.MustCompile(textP)
	}
	if s, ok := d.(string); ok {
		return textRe.Match([]byte(s))
	}
	return false
}

func (s *SchemaValidator) validate_ip(d interface{}) bool {
	if s, ok := d.(string); ok {
		return net.ParseIP(s) != nil
	}
	return false
}

func (s *SchemaValidator) validate_array(a []interface{}, attr *SchemaItem, r string) (bool, CcError) {
	if attr.Required && len(a) == 0 {
		return false, nil
	}
	for _, ele := range a {
		if valid, ce := s.validate_data(ele, attr, r); !valid {
			return false, ce
		}
	}
	return true, nil
}

func (s *SchemaValidator) validate_data(p interface{}, attr *SchemaItem, r string) (bool, CcError) {
	var valid bool
	var ce CcError = nil
	switch attr.Type {
	case T_Int:
		valid = s.validate_int(p)
	case T_String:
		valid = s.validate_string(p) && (!attr.Required || p.(string) != "")
	case T_Uuid:
		valid = s.validate_uuid(p)
	case T_UuidNull:
		valid = s.validate_uuidNull(p)
	case T_Text:
		valid = s.validate_text(p) && (!attr.Required || p.(string) != "")
	case T_IP:
		valid = s.validate_ip(p)
	case T_Bool:
		valid = validate_bool(p)
	case T_Struct:
		if d, ok := p.(map[string]interface{}); ok {
			valid, ce = s.validate_struct(d, attr.SubItem.(map[string]*SchemaItem), r)
		} else {
			valid = false
		}
	case T_Filter:
		if fs, ok := p.([]interface{}); ok {
			valid, ce = s.validate_filter(fs, attr.SubItem.(map[string]*SchemaItem), r)
		} else {
			valid = false
			ce = NewError(r, ErrInvalid, ErrFilter, ErrNull)
		}
	case T_Order:
		if fs, ok := p.([]interface{}); ok {
			valid, ce = s.validate_order(fs, attr.SubItem.(map[string]*SchemaItem), r)
		} else {
			valid = false
			ce = NewError(r, ErrInvalid, ErrOrder, ErrNull)
		}

	case T_Array:
		if a, ok := p.([]interface{}); ok {
			valid, ce = s.validate_array(a, attr.SubItem.(*SchemaItem), r)
		} else {
			valid = false
		}
	default:
		panic("Parameter Format Invalid")
	}
	if valid {
		valid = (attr.Val == nil) || attr.Val(p)
	}
	return valid, ce
}

func (s *SchemaValidator) validate_struct(m map[string]interface{}, attrs map[string]*SchemaItem, r string) (bool, CcError) {
	for k, v := range attrs {
		if p, ok := m[k]; ok && p != nil {
			if valid, ce := s.validate_data(p, v, k); !valid {
				if ce == nil {
					ce = NewError(r, ErrInvalid, k, ErrNull)
				}
				return false, ce
			}
		} else {
			if v.Required {
				return false, NewError(r, ErrMiss, k, ErrNull)
			}
		}
	}
	return true, nil
}

func (s *SchemaValidator) validate_filter(fs []interface{}, attrs map[string]*SchemaItem, r string) (bool, CcError) {
	return s.validate_fo(fs, attrs, r, "field", "value")
}

func (s *SchemaValidator) validate_order(fs []interface{}, attrs map[string]*SchemaItem, r string) (bool, CcError) {
	return s.validate_fo(fs, attrs, r, "field", "direction")
}

func (s *SchemaValidator) validate_fo(fs []interface{}, attrs map[string]*SchemaItem, r string, key string, value string) (bool, CcError) {
	for _, f := range fs {
		if fm, ok := f.(map[string]interface{}); !ok {
			return false, NewError(r, ErrInvalid, ErrFilter, ErrNull)
		} else {
			field, ok := fm[key]
			if !ok {
				return false, nil
			}
			value, ok := fm[value]
			if !ok {
				return false, nil
			}
			subattr, ok := attrs[field.(string)]
			if !ok {
				return false, nil
			}
			if valid, ce := s.validate_data(value, subattr, r); !valid {
				return valid, ce
			}
		}
	}
	return true, nil
}

func (s *SchemaValidator) Validate(attrs map[string]*SchemaItem, rr io.Reader, r string, d interface{}) CcError {
	bt, err := ioutil.ReadAll(rr)
	if err != nil {
		return NewError(r, ErrMalformed, ErrNull, err.Error())
	}
	var f interface{}
	if err := json.Unmarshal(bt, &f); err != nil {
		return NewError(r, ErrMalformed, ErrNull, err.Error())
	}
	m := f.(map[string]interface{})
	if _, ce := s.validate_struct(m, attrs, r); ce != nil {
		return ce
	}
	if err := json.Unmarshal(bt, d); err != nil {
		return NewError(r, ErrMalformed, ErrNull, err.Error())
	}
	return nil
}

// Entry
func Validate(schema map[string]*SchemaItem, rr io.Reader, r string, d interface{}) CcError {
	if schemaVal == nil {
		return NewSysErr(errors.New("Schema validator has not been init."))
	}
	return schemaVal.Validate(schema, rr, r, d)
}

// Validators
func ValIntRange(s, e int) ValidatorFunc {
	return func(d interface{}) bool {
		df, ok := d.(float64)
		if !ok {
			return false
		}
		di := int(df)
		return di >= s && di <= e
	}
}

func ValIntEles(eles ...int) ValidatorFunc {
	em := make(map[int]bool)
	for _, ele := range eles {
		em[ele] = true
	}
	return func(d interface{}) bool {
		df, ok := d.(float64)
		if !ok {
			return false
		}
		di := int(df)
		_, val := em[di]
		return val
	}
}

func ValRegExp(ptn string) ValidatorFunc {
	return func(d interface{}) bool {
		if s, ok := d.(string); ok {
			matched, err := regexp.MatchString(ptn, s)
			return err == nil && matched
		}
		return false
	}
}

func ValStrLen(le int) ValidatorFunc {
	return func(d interface{}) bool {
		if s, ok := d.(string); ok {
			return len(s) <= le
		}
		return false
	}
}

func ValTextLen(le int) ValidatorFunc {
	return func(d interface{}) bool {
		if s, ok := d.(string); ok {
			return utf8.RuneCountInString(s) <= le
		}
		return false
	}
}

func ValStrEq(s string) ValidatorFunc {
	return func(d interface{}) bool {
		if sd, ok := d.(string); ok {
			return sd == s
		}
		return false
	}
}

func MultiValid(valids ...Validater) ValidatorFunc {
	return func(d interface{}) bool {
		for _, valid := range valids {
			if !valid(d) {
				return false
			}
		}
		return true
	}
}
