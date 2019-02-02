package uuid

import (
	"testing"
)

const (
	nameSpace = "test"
	key       = "test"
)

func Test_GenerateUUID5(t *testing.T) {
	e := GenerateUUID5(nameSpace, key)
	e.Base62()
}

func Benchmark_GenerateUUID5(b *testing.B) {
	e := GenerateUUID5(nameSpace, key)
	for i := 0; i < b.N; i++ { //use b.N for looping
		e.Base62()
	}
}

func Test_GenerateUUID4(t *testing.T) {
	e := GenerateUUID4()
	e.Base62()
}

func Benchmark_GenerateUUID4(b *testing.B) {
	e := GenerateUUID4()
	for i := 0; i < b.N; i++ { //use b.N for looping
		e.Base62()
	}
}
