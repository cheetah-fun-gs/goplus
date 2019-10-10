package uuid

import (
	"testing"
)

const (
	nameSpace = "test"
	key       = "test"
)

func Test_NewV5(t *testing.T) {
	e := NewV5(nameSpace, key)
	e.Base62()
}

func Benchmark_NewV5(b *testing.B) {
	e := NewV5(nameSpace, key)
	for i := 0; i < b.N; i++ { //use b.N for looping
		e.Base62()
	}
}

func Test_NewV4(t *testing.T) {
	e := NewV4()
	e.Base62()
}

func Benchmark_NewV4(b *testing.B) {
	e := NewV4()
	for i := 0; i < b.N; i++ { //use b.N for looping
		e.Base62()
	}
}
