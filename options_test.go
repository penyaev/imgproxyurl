package imgproxyurl

import "testing"

type ts struct{}

func (ts) String() string { return "tstring" }

func Test_format(t *testing.T) {
	type args struct {
		key       string
		arguments []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "no arguments", args: args{key: "k", arguments: nil}, want: ""},
		{name: "one int argument", args: args{key: "k", arguments: []interface{}{1}}, want: "1"},
		{name: "mixed-type arguments", args: args{key: "k", arguments: []interface{}{1, "z", ts{}}}, want: "1:z:tstring"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.args.key, tt.args.arguments...); got != tt.want {
				t.Errorf("format() = %v, want %v", got, tt.want)
			}
		})
	}
}
