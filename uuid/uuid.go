package uuid

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/google/uuid"
	"gitlab.liebaopay.com/mikezhang/goplus/encoding/basen"
)

func calMD5(value string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(value))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// UUID UUID
type UUID []byte

// Base64 to Base64
func (u *UUID) Base64() string {
	return base64.RawURLEncoding.EncodeToString(*u)
}

// Base62 to Base62
func (u *UUID) Base62() string {
	e := basen.NewBase62()
	return e.Encode(*u)
}

// MD5 to MD5
func (u *UUID) MD5() string {
	return hex.EncodeToString(*u)
}

// GenerateUUID5 生成 UUID version 5
func GenerateUUID5(nameSpace string, token string) UUID {
	nameSpaceUUID, err := uuid.Parse(calMD5(nameSpace))
	if err != nil {
		panic(err)
	}
	u := uuid.NewSHA1(nameSpaceUUID, []byte(token))
	b, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return UUID(b)
}

// GenerateUUID4 生成 UUID version 4
func GenerateUUID4() UUID {
	u := uuid.New()
	b, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return UUID(b)
}
