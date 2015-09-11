package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"unicode/utf16"
)

func HexToInt(b []byte) uint64 {

	var result uint64

	if len(b) == 0 {
		return result
	}

	for _, symbol := range b {
		if symbol >= '0' && symbol <= '9' {
			result <<= 4
			result += uint64(symbol - '0')
		} else if symbol >= 'a' && symbol <= 'f' {
			result <<= 4
			result += uint64(symbol - 'a' + 10)
		} else {
			break
		}
	}

	return result
}

func BytesToInt64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}

func Int64ToBytes(data int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)
	return buf.Bytes()
}

func BytesToInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}

func Int32ToBytes(val int32) []byte {
	return int32ToBytesWithRoute(val, binary.LittleEndian)
}

func Int32ToHexToBytes(val int32) []byte {

	b := int32ToBytesWithRoute(val, binary.BigEndian)
	h := hex.EncodeToString(b)

	res := make([]byte, len(h))
	copy(res, h)

	return res
}

func int32ToBytesWithRoute(data int32, route binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, route, data)
	return buf.Bytes()
}

func Utf16ToUtf8(data []byte) string {

	symbols := int(len(data) / 2)
	u16 := make([]uint16, symbols)

	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &u16)

	return string(utf16.Decode(u16))
}

func Utf8ToUtf16(u8 []rune) []byte {

	u16 := utf16.Encode(u8)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, u16)
	return buf.Bytes()
}

func IntToString(val int) string {
	return fmt.Sprintf("%d", val)
}
