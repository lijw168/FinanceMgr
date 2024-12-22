package util

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GBKToUTF8(gbkBuf []byte) ([]byte, error) {
	sReader := bytes.NewReader(gbkBuf)
	dReader := transform.NewReader(sReader, simplifiedchinese.GBK.NewDecoder())
	utfBuf, err := io.ReadAll(dReader)
	if err != nil {
		return nil, err
	}
	return utfBuf, nil
}

func UTF8ToGBK(utfBuf []byte) ([]byte, error) {
	sReader := bytes.NewReader(utfBuf)
	dReader := transform.NewReader(sReader, simplifiedchinese.GBK.NewEncoder())
	gbkBuf, err := io.ReadAll(dReader)
	if err != nil {
		return nil, err
	}
	return gbkBuf, nil
}
