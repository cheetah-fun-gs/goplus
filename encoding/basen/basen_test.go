package basen

import (
	"testing"
)

var data = []byte{12, 34, 56, 78, 90, 13, 14, 15, 16, 17, 18, 21, 32, 45, 67, 89}

func Test_Base62(t *testing.T) {
	e := NewBase62()
	e.Encode(data)
}

func Benchmark_Base62(b *testing.B) {
	e := NewBase62()
	for i := 0; i < b.N; i++ { //use b.N for looping
		e.Encode(data)
	}
}
