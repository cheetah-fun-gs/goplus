package reflect

import (
	"reflect"
	"testing"
	"time"
)

func TestMockKeep(t *testing.T) {
	type test0 struct {
		A0 int    `json:"a0,omitempty"`
		B0 string `json:"b0,omitempty"`
	}
	type test1 struct {
		A  int    `json:"a,omitempty"`
		B  string `json:"b,omitempty"`
		T0 test0
	}
	type test2 struct {
		A1 *test1
		B1 map[string]int
		C  [2]int
		D  []time.Time
	}

	testWant := map[string]interface{}{
		"A1": map[string]interface{}{
			"a": 0,
			"b": "",
			"T0": map[string]interface{}{
				"a0": 0,
				"b0": "",
			},
		},
		"B1": map[string]interface{}{
			"": 0,
		},
		"C": []interface{}{0, 0},
		"D": []interface{}{time.Time{}},
	}

	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "basic",
			args: args{
				v: 1,
			},
			want: 1,
		},
		{
			name: "all",
			args: args{
				v: &test2{},
			},
			want: testWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MockKeep(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MockKeep() = %v, want %v", got, tt.want)
			}
		})
	}
}
