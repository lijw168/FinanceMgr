package encoding

import (
	"encoding/binary"
	"financeMgr/src/common/utils"
)

func MarshalBool(buf []byte, b bool) int {
	if b {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return 1
}

func MarshalUInt8(buf []byte, v uint8) int {
	buf[0] = v
	return 1
}

func MarshalUInt32(buf []byte, v uint32) int {
	binary.LittleEndian.PutUint32(buf, v)
	return 4
}

func MarshalUInt64(buf []byte, v uint64) int {
	binary.LittleEndian.PutUint64(buf, uint64(v))
	return 8
}

func UnmarshalBool(buf []byte, b *bool) int {
	if buf[0] == 1 {
		*b = true
	} else if buf[0] == 0 {
		*b = false
	} else {
		utils.Assert(false)
	}
	return 1
}

func UnmarshalUInt8(buf []byte, v *uint8) int {
	*v = buf[0]
	return 1
}

func UnmarshalUInt32(buf []byte, v *uint32) int {
	*v = binary.LittleEndian.Uint32(buf)
	return 4
}

func UnmarshalUInt64(buf []byte, v *uint64) int {
	*v = uint64(binary.LittleEndian.Uint64(buf))
	return 8
}

func UnmarshalInt8(buf []byte, v *int8) int {
	*v = int8(buf[0])
	return 1
}

func UnmarshalInt32(buf []byte, v *int32) int {
	*v = int32(binary.LittleEndian.Uint32(buf))
	return 4
}

func UnmarshalInt64(buf []byte, v *int64) int {
	*v = int64(binary.LittleEndian.Uint64(buf))
	return 8
}

func MarshalStr(buf []byte, v string, size int) int {
	for i := 0; i < size; i++ {
		buf[i] = 0
	}
	copy(buf, v)
	return size
}

func UnmarshalStr(buf []byte, v *string, size int) int {
	i := 0
	for ; i < size; i++ {
		if buf[i] == 0 {
			break
		}
	}

	*v = string(buf[:i])
	return size
}
