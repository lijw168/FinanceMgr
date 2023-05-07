package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"financeMgr/src/analysis-server/cli/tbase"
	"fmt"
	"os"
	"reflect"

	//"sort"
	"time"
)

func FormatViewOutput(v interface{}) {
	if v == nil {
		return
	}
	bt, _ := json.Marshal(v)
	m := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(bt))
	decoder.UseNumber()
	decoder.Decode(&m)
	table := tbase.NewWriter(os.Stdout)

	table.SetHeader([]string{"Field", "Value"})
	table.SetAutoWrapText(false)

	header := make([]string, 0, 5)
	for k := range m {
		header = append(header, k)
	}
	//sort.Strings(header)
	for _, h := range header {
		if m[h] == nil {
			fmt.Printf("please check program code,the h is %v, m[h] is nil\n", h)
			continue
		}
		switch reflect.TypeOf(m[h]).Kind() {
		case reflect.Slice, reflect.Map:
			bt, _ = json.Marshal(m[h])
			table.Append([]string{h, string(bt)})
		default:
			table.Append([]string{h, fmt.Sprintf("%v", m[h])})
		}
	}
	table.Render()
}

func FormatStringListOutput(field string, data []string) {
	table := tbase.NewWriter(os.Stdout)
	table.SetHeader([]string{field})
	table.SetAutoWrapText(false)

	for _, v := range data {
		table.Append([]string{v})
	}
	table.Render()
}

func FormatListOutput(header []string, l interface{}) {
	if l == nil {
		return
	}
	ls := reflect.ValueOf(l)
	if ls.Kind() != reflect.Slice {
		return
	}
	header = mapString(header, unix2Camel)
	// Validate header exist or not
	validHdr, err := validateField(l, header)
	if err != nil {
		FormatErrorOutput(err)
		return
	}
	if len(validHdr) > 0 {
		table := tbase.NewWriter(os.Stdout)
		headerLen := len(validHdr)
		for i := 0; i < ls.Len(); i++ {
			row := make([]string, 0, headerLen)
			for _, h := range validHdr {
				v := reflect.Indirect(ls.Index(i)).FieldByName(h)
				switch v.Kind() {
				case reflect.Slice, reflect.Map:
					bt, _ := json.Marshal(v.Interface())
					row = append(row, string(bt))
				case reflect.Ptr:
					bt, _ := json.Marshal(reflect.Indirect(v).Interface())
					row = append(row, string(bt))
				case reflect.Struct:
					if v.Type().String() == "time.Time" {
						row = append(row, FormatTime(v.Interface().(time.Time)))
					} else {
						bt, _ := json.Marshal(v.Interface())
						row = append(row, string(bt))

					}
				default:
					row = append(row, fmt.Sprintf("%v", v))
				}
			}
			table.Append(row)
		}
		table.SetHeader(validHdr)
		table.Render()
	}
	// Print invalid header warning
	hdrSet := make(map[string]bool)
	for _, h := range header {
		hdrSet[h] = false
	}
	for _, h := range validHdr {
		hdrSet[h] = true
	}
	for h, valid := range hdrSet {
		if !valid {
			FormatMessageOutput("Warn: Unsupport column '" + h + "'")
		}
	}
}

func FormatErrorOutput(err error) {
	fmt.Println(err.Error())
}

func FormatMessageOutput(msg string) {
	fmt.Println(msg)
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func unix2Camel(s string) string {
	var tmp = &bytes.Buffer{}
	var upper bool = true
	for _, c := range s {
		if c == '-' || c == '_' {
			upper = true
		} else if upper {
			if c >= 'a' && c <= 'z' {
				c -= 32
			}
			tmp.WriteRune(c)
			upper = false
		} else {
			tmp.WriteRune(c)
		}
	}
	return tmp.String()
}

func mapString(in []string, h func(string) string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		out = append(out, h(s))
	}
	return out
}

func validateField(i interface{}, fields []string) ([]string, error) {
	validF := make([]string, 0, len(fields))
	tp := reflect.TypeOf(i)
	for {
		switch tp.Kind() {
		case reflect.Slice:
			tp = tp.Elem()
		case reflect.Ptr:
			tp = tp.Elem()
		case reflect.Struct:
			goto end
		default:
			return nil, errors.New("Bad type to validate field")
		}
	}
end:
	for _, f := range fields {
		if _, valid := tp.FieldByName(f); valid {
			validF = append(validF, f)
		}
	}
	return validF, nil
}
