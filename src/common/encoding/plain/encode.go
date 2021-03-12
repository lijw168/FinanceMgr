package plain

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

/*
Package plain encodes instance in plain k-v format.
*/

const INDENT_STRING = "  "
const NULL_STRING = "null"

func Marshal(v interface{}) ([]byte, error) {
	return newEncodeState().marshal(v)
}

type encoder struct {
	bytes.Buffer
}

func newEncodeState() *encoder {
	return &encoder{}
}

func (e *encoder) marshal(i interface{}) ([]byte, error) {
	val := reflect.ValueOf(i)
	val = reflect.Indirect(val)

	// Write the overall title
	e.writeStructTypeLine(val.Type(), "")

	e.doMarshal("", val, "")
	return e.Bytes(), nil
}

func (e *encoder) writeStructTypeLine(t reflect.Type, indent string) {
	e.WriteString(indent)
	e.WriteString(t.Name())
	e.WriteString(":\n")
}

func (e *encoder) doMarshal(fieldName string, val reflect.Value, indent string) {
	switch val.Kind() {
	case reflect.Bool:
		valStr := e.encodeBool(val)
		e.writeFieldValueLine(indent, fieldName, valStr)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		valStr := e.encodeInt(val)
		e.writeFieldValueLine(indent, fieldName, valStr)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		valStr := e.encodeUint(val)
		e.writeFieldValueLine(indent, fieldName, valStr)

	case reflect.String:
		valStr := e.encodeString(val)
		e.writeFieldValueLine(indent, fieldName, valStr)

	case reflect.Ptr:
		if val.IsNil() {
			e.writeFieldValueLine(indent, fieldName, NULL_STRING)
			return
		}

		e.doMarshal(fieldName, val.Elem(), indent)

	case reflect.Slice:
		if val.IsNil() {
			e.writeFieldValueLine(indent, fieldName, NULL_STRING)
			return
		}

		e.WriteString(indent)
		e.writeStructFieldName(fieldName)

		for i := 0; i < val.Len(); i++ {
			fieldVal := val.Index(i)
			fieldIndex := strconv.FormatInt(int64(i), 10)
			e.doMarshal(fieldIndex, fieldVal, indent+INDENT_STRING)
		}

	case reflect.Struct:
		e.marshalStruct(fieldName, val, indent)
	}
}

func (e *encoder) marshalStruct(fieldName string, val reflect.Value, indent string) {
	e.WriteString(indent)
	e.writeStructFieldName(fieldName)

	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := t.Field(i).Name

		if e.isExportedFiled(fieldName) {
			fieldTag := t.Field(i).Tag
			fieldVal := val.Field(i)
			fieldIndent := indent + INDENT_STRING

			tagKV := e.parseFieldTag(fieldTag)
			e.preMarshalField(tagKV, fieldIndent)
			e.doMarshal(fieldName, fieldVal, fieldIndent)
			e.postMarshalField(tagKV, fieldIndent)
		}
	}
}

const (
	plain_title   = "title"
	plain_newline = "newline"
)

func (e *encoder) parseFieldTag(tag reflect.StructTag) map[string]string {
	m := make(map[string]string, 0)

	plainTag := tag.Get("plain")
	if plainTag == "" {
		return m
	}

	tags := strings.Split(plainTag, ",")
	for _, tag := range tags {
		kv := strings.Split(tag, "=")
		if len(kv) == 1 {
			// single tag which has no value.
			k := strings.TrimSpace(kv[0])
			m[k] = ""
		} else {
			k := strings.TrimSpace(kv[0])
			v := strings.TrimSpace(kv[1])
			m[k] = v
		}
	}
	return m
}

func (e *encoder) preMarshalField(kvs map[string]string, indent string) {
	if _, ok := kvs[plain_title]; ok {
		e.WriteString(indent)
		e.WriteString(kvs[plain_title])
		e.newLine()
	}
}

func (e *encoder) postMarshalField(kvs map[string]string, indent string) {
	if _, ok := kvs[plain_newline]; ok {
		e.newLine()
	}
}

func (e *encoder) isExportedFiled(fieldName string) bool {
	firstLetter := fieldName[0]
	if firstLetter >= 'A' && firstLetter <= 'Z' {
		return true
	}
	return false
}

func (e *encoder) writeFieldValueLine(indent, field, value string) {
	e.WriteString(indent)
	e.writeFieldName(field)
	e.WriteString(value)
	e.newLine()
}

func (e *encoder) writeFieldName(fieldName string) {
	e.WriteString(fieldName)
	e.WriteString(": ")
}

func (e *encoder) writeStructFieldName(fieldName string) {
	if fieldName != "" {
		e.WriteString(fieldName)
		e.WriteByte(':')
		e.newLine()
	}
}

func (e *encoder) newLine() {
	e.WriteByte('\n')
}

func (e *encoder) encodeBool(v reflect.Value) string {
	if v.Bool() {
		return "true"
	} else {
		return "false"
	}
}

func (e *encoder) encodeInt(v reflect.Value) string {
	return strconv.FormatInt(v.Int(), 10)
}

func (e *encoder) encodeUint(v reflect.Value) string {
	return strconv.FormatUint(v.Uint(), 10)
}

func (e *encoder) encodeString(v reflect.Value) string {
	return v.String()
}
