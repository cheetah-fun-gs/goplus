package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func fillByte(in []byte, b byte, n int) []byte {
	out := []byte{}
	for i := 0; i < n; i++ {
		out = append(out, b)
	}
	for _, b := range in {
		out = append(out, b)
	}
	return out
}

// ByteToUint64 字节转整型
func ByteToUint64(in []byte) (uint64, error) {
	if len(in) > 8 {
		return 0, fmt.Errorf("in length more than 8")
	}
	if len(in) < 8 {
		in = fillByte(in, 0x00, 8-len(in))
	}
	buffer := bytes.NewBuffer(in)
	var out uint64
	err := binary.Read(buffer, binary.BigEndian, &out)
	if err != nil {
		return 0, err
	}
	return out, nil
}

// Uint64ToByte 整型转字节
func Uint64ToByte(i uint64) ([]byte, error) {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
