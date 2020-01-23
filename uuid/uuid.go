// Package uuid google uuid更方便的二次封装, 含uuid4和uuid5两种算法, 以及md5, base62, base64三种格式
package uuid

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/nicksnyder/basen"
)

// md5 to xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func md5ToUUIDString(md5Str string) string {
	uuidSplit := []string{}
	md5Split := strings.Split(md5Str, "")
	for i, char := range md5Split {
		if i == 8 || i == 12 || i == 16 || i == 20 {
			uuidSplit = append(uuidSplit, "-")
		}
		uuidSplit = append(uuidSplit, char)
	}
	return strings.Join(uuidSplit, "")
}

func calMD5(value string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(value))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// UUID UUID
type UUID []byte

// Base64 to Base64
func (u UUID) Base64() string {
	return base64.RawURLEncoding.EncodeToString(u)
}

// Base62 to Base62
func (u UUID) Base62() string {
	return basen.Base62Encoding.EncodeToString(u)
}

// MD5 to MD5
func (u UUID) MD5() string {
	return hex.EncodeToString(u)
}

// String to form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], u)
	return string(buf[:])
}

// NewV5 生成 UUID version 5
func NewV5(nameSpace string, key string) UUID {
	nameSpaceUUID, err := uuid.Parse(md5ToUUIDString(calMD5(nameSpace)))
	if err != nil {
		panic(err)
	}
	u := uuid.NewSHA1(nameSpaceUUID, []byte(key))
	b, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return UUID(b)
}

// NewV4 生成 UUID version 4
func NewV4() UUID {
	u := uuid.New()
	b, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return UUID(b)
}
