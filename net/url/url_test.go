package url

import "testing"

func TestToRawQuery(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "map[string]interface{}",
			args: args{
				v: map[string]interface{}{
					"a": 1,
					"b": "2",
					"c": []string{"3", "4"},
					"d": []int{5, 6},
				},
			},
			want:    "a=1&b=2&c=3&c=4&d=5&d=6",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToRawQuery(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToRawQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToRawQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
