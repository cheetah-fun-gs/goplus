package binary

import (
	"bytes"
	"encoding/binary"
)

// ByteFillLeft 左边填充 byte
func ByteFillLeft(in []byte, b byte, n int) []byte {
	out := []byte{}
	for i := 0; i < n; i++ {
		out = append(out, b)
	}
	for _, b := range in {
		out = append(out, b)
	}
	return out
}

// ByteFillRight 右边填充 byte
func ByteFillRight(in []byte, b byte, n int) []byte {
	for i := 0; i < n; i++ {
		in = append(in, b)
	}
	return in
}

// ByteToUint64 字节转整型
func ByteToUint64(in []byte) uint64 {
	if len(in) > 8 {
		panic("in length more than 8")
	}
	if len(in) < 8 {
		in = ByteFillLeft(in, 0x00, 8-len(in))
	}
	buffer := bytes.NewBuffer(in)
	var out uint64
	err := binary.Read(buffer, binary.BigEndian, &out)
	if err != nil {
		panic(err)
	}
	return out
}

// Uint64ToByte 整型转字节
func Uint64ToByte(i uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, i)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}
