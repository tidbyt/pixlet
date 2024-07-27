package members

import (
	"reflect"
	"testing"
)

func TestParseParameters(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]string
	}{
		{"Test1", args{p: "HELLO;KEY1=value1;KEY2=value2"}, "HELLO", map[string]string{"KEY1": "value1", "KEY2": "value2"}},
		{"Test2", args{p: "TEST2;GAS=FUEL;ROCK=STONE"}, "TEST2", map[string]string{"GAS": "FUEL", "ROCK": "STONE"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseParameters(tt.args.p)
			if got != tt.want {
				t.Errorf("ParseParameters() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseParameters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParseRecurrenceParams(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]string
	}{
		{"Test1", args{p: "RRULE:FREQ=WEEKLY;INTERVAL=2;WKST=SU;COUNT=10"}, "RRULE", map[string]string{"FREQ": "WEEKLY", "COUNT": "10", "WKST": "SU", "INTERVAL": "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseRecurrenceParams(tt.args.p)
			if got != tt.want {
				t.Errorf("ParseRecurrenceParams() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseRecurrenceParams() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUnescapeString(t *testing.T) {
	type args struct {
		l string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Test1", args{l: `Hello\, world\; lorem \\ipsum.`}, `Hello, world; lorem \ipsum.`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnescapeString(tt.args.l); got != tt.want {
				t.Errorf("UnescapeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
