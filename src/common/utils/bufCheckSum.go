package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"common/constant"
)

func BufCheckSum(buf []byte, method int) string {
	switch method {
	case constant.BUF_CHECK_SUM_CRC:
		bufCrc := crc32.ChecksumIEEE(buf)
		return fmt.Sprintf("%d", bufCrc)
	case constant.BUF_CHECK_SUM_MD5:
		bufMd5 := md5.Sum(buf)
		return hex.EncodeToString(bufMd5[:])
	default:
		Assert(false)
	}

	return ""
}
