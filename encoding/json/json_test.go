package json

import (
	"testing"
)

func TestStringsToMap(t *testing.T) {
	type args struct {
		datas []string
		v     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "StringsToMap1",
			args: args{
				datas: []string{"a", "\"b\"", "c", "\"d\""},
				v:     &map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "StringsToMap2",
			args: args{
				datas: []string{"a", "1", "c", "2"},
				v:     &map[string]int{},
			},
			wantErr: false,
		},
		{
			name: "StringsToMap3",
			args: args{
				datas: []string{"a", "1", "c", "\"2\""},
				v:     &map[string]interface{}{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StringsToMap(tt.args.datas, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("StringsToMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringsToList(t *testing.T) {
	type args struct {
		datas []string
		v     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "StringsToList1",
			args: args{
				datas: []string{"\"a\"", "1", "\"c\"", "2"},
				v:     &[]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "StringsToList2",
			args: args{
				datas: []string{"\"a\"", "\"1\"", "\"c\"", "\"2\""},
				v:     &[]string{},
			},
			wantErr: false,
		},
		{
			name: "StringsToList3",
			args: args{
				datas: []string{"1", "1", "3", "4"},
				v:     &[]int{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StringsToList(tt.args.datas, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("StringsToList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
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
			name: "ToJSON1",
			args: args{
				map[interface{}]interface{}{
					"1": "1",
					2:   2,
				},
			},
			want:    "{\"1\":\"1\",\"2\":2}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSON(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	type args struct {
		data string
		v    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FromJSONList",
			args: args{
				data: "[1,2,3,4]",
				v:    &[]int{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FromJSON(tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FromJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
