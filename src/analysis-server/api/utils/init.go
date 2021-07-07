package utils

import (
	"regexp"
)

func init() {
	strParaRe = regexp.MustCompile(strParaPattern)
	//hnParaRe = regexp.MustCompile(hnParaPattern)
	commonIdParaRe = regexp.MustCompile(CommonIdParaPattern)
}
