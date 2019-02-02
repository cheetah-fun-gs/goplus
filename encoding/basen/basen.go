// Package basen 将任意数据按指定的字符集编码成可见字符
// 算法思路：
// 1. 根据字符集的数量，算出每8个字节需要多少个字符表示
// 2. 余数的值加在最后一位
package basen

import (
	"fmt"
	"math"
	"strings"

	"gitlab.liebaopay.com/mikezhang/goplus/encoding/binary"
	"gitlab.liebaopay.com/mikezhang/goplus/number"

	stringsplus "gitlab.liebaopay.com/mikezhang/goplus/strings"
)

// 字符集
const (
	CharsetBase62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// Encoder 编码器
type Encoder struct {
	Charset   []string       // 字符集
	CharIndex map[string]int // 字符索引
	CharNum   int            // 每8字节所需字符数
}

func (e *Encoder) decimalToAny(num uint64) string {
	n := uint64(len(e.Charset))
	newNumStr := ""
	var remainder uint64
	var remainderString string
	for num != 0 {
		remainder = num % n
		remainderString = e.Charset[remainder]
		newNumStr = remainderString + newNumStr
		num = num / n
	}
	return newNumStr
}

func (e *Encoder) anyToDecimal(any string) (uint64, error) {
	n := len(e.Charset)
	num := 0
	for index, value := range strings.Split(any, "") {
		valueNum, ok := e.CharIndex[value]
		if !ok {
			return 0, fmt.Errorf("%s is not in charset", value)
		}
		num = num + valueNum*number.Pow(n, len(any)-1-index)
	}
	return uint64(num), nil
}

func (e *Encoder) calCharNum() int {
	var charNum int
	var int64Max = math.Pow(256.0, 8.0)
	for {
		if math.Pow(float64(len(e.Charset)), float64(charNum)) >= int64Max {
			break
		}
		charNum++
	}
	return charNum
}

func (e *Encoder) init() {
	e.CharIndex = make(map[string]int)
	for index, value := range e.Charset {
		e.CharIndex[value] = index
	}
	e.CharNum = e.calCharNum()
}

// NewEncoder 一个新的编码器
func NewEncoder(charset string) (*Encoder, error) {
	if len(charset) < 8 {
		return nil, fmt.Errorf("charset length is less than 8")
	}
	e := Encoder{
		Charset: strings.Split(charset, ""),
	}
	e.init()
	return &e, nil
}

// Encode 编码
func (e *Encoder) Encode(in []byte) string {
	out := []string{}
	offset := 0
	remainder := len(in) % 8
	for {
		if offset >= len(in) {
			break
		}
		var b []byte
		var last bool
		if offset+8 <= len(in) {
			b = in[offset : offset+8]
		} else {
			b = in[offset:]
			last = true
		}
		s := e.decimalToAny(binary.ByteToUint64(b))
		if len(s) < e.CharNum && !last {
			s = stringsplus.StringFillLeft(s, e.Charset[0], e.CharNum-len(s))
		} else if len(s) > e.CharNum {
			panic("decimalToAny length more than CharNum")
		}
		out = append(out, s)
		offset += 8
	}
	remainderString := e.Charset[remainder]
	out = append(out, remainderString)
	return strings.Join(out, "")
}

// Decode 解码
func (e *Encoder) Decode(in string) ([]byte, error) {
	inSplit := strings.Split(in, "")
	out := []byte{}
	offset := 0
	remainderString := inSplit[len(inSplit)-1]
	remainder, ok := e.CharIndex[remainderString]
	if !ok {
		return nil, fmt.Errorf("%s is not in charset", remainderString)
	}
	inSplit = inSplit[:len(inSplit)-1]
	for {
		if offset >= len(inSplit) {
			break
		}
		var sSplit []string
		var last bool
		if offset+e.CharNum <= len(inSplit) {
			sSplit = inSplit[offset : offset+e.CharNum]
		} else {
			sSplit = inSplit[offset:]
			last = true
		}

		bUnit64, err := e.anyToDecimal(strings.Join(sSplit, ""))
		if err != nil {
			return nil, err
		}
		b := binary.Uint64ToByte(bUnit64)
		if last {
			b = b[len(b)-remainder:]
		}
		for _, bb := range b {
			out = append(out, bb)
		}
		offset += e.CharNum
	}
	return out, nil
}

// NewBase62 一个新的编码器
func NewBase62() *Encoder {
	e, _ := NewEncoder(CharsetBase62)
	return e
}
