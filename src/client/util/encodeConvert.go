package util

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func GBKToUTF8(gbkBuf []byte) ([]byte, error) {
	sReader := bytes.NewReader(gbkBuf)
	dReader := transform.NewReader(sReader, simplifiedchinese.GBK.NewDecoder())
	utfBuf, err := ioutil.ReadAll(dReader)
	if err != nil {
		return nil, err
	}
	return utfBuf, nil
}

func UTF8ToGBK(utfBuf []byte) ([]byte, error) {
	sReader := bytes.NewReader(utfBuf)
	dReader := transform.NewReader(sReader, simplifiedchinese.GBK.NewEncoder())
	gbkBuf, err := ioutil.ReadAll(dReader)
	if err != nil {
		return nil, err
	}
	return gbkBuf, nil
}
