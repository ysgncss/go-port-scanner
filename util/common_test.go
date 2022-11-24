package util

import (
	"testing"
)

func TestValidIPV4Address(t *testing.T) {
	type args struct {
		queryIP string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test",
			args: args{queryIP: "2.0.0.111"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidIpv4(tt.args.queryIP); got != tt.want {
				t.Errorf("ValidIPV4Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCidrIpRange(t *testing.T) {
	type args struct {
		cidr string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{cidr: "192.168.31.1/24"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCidrIpRange(tt.args.cidr)
			t.Logf("GetCidrIpRange() = %v", got)

		})
	}
}

func TestPing(t *testing.T) {
	type args struct {
		target  string
		timeout int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{target: "www.baidu.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ping(tt.args.target)
			t.Logf("GetCidrIpRange() = %v", got)

		})
	}
}
