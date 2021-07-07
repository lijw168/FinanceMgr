package utils

import (
	"regexp"
)

const strParaPattern string = `^([a-zA-Z0-9\_-]|[\p{Han}])*$`
const CommonIdParaPattern string = `^([0-9])*$`

//const hnParaPattern string = `^([a-zA-Z0-9\_\-\.])*$`

var strParaRe *regexp.Regexp

//var hnParaRe *regexp.Regexp
var commonIdParaRe *regexp.Regexp

func VerStrP(s string) bool {
	return strParaRe.Match([]byte(s))
}

func VerCommonIdP(s string) bool {
	return commonIdParaRe.Match([]byte(s))
}
