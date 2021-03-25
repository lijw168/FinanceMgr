package utils

import (
	"regexp"
)

const strParaPattern string = `^([a-zA-Z0-9\_-]|[\p{Han}])*$`
const hnParaPattern string = `^([a-zA-Z0-9\_\-\.])*$`

var strParaRe *regexp.Regexp
var hnParaRe *regexp.Regexp

func VerStrP(s string) bool {
	return strParaRe.Match([]byte(s))
}

// func VerHostnameP(s string) bool {
// 	return hnParaRe.Match([]byte(s))
// }
