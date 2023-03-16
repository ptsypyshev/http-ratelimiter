package app

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_getIP(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Standard IP (X-Forwarded-For)",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Header: http.Header{
							"X-Forwarded-For": []string{"10.20.30.40"},
						},
					},
				},
			},
			want: "10.20.30.40",
		},
		{
			name: "Standard IP (Remote IP)",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						RemoteAddr: "192.168.1.10:8080",
					},
				},
			},
			want: "192.168.1.10",
		},
		{
			name: "Both IP (Remote IP and X-Forwarded-For)",
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						RemoteAddr: "192.168.1.10:8080",
						Header: http.Header{
							"X-Forwarded-For": []string{"10.20.30.40"},
						},
					},					
				},
			},
			want: "10.20.30.40",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIP(tt.args.c); got != tt.want {
				t.Errorf("getIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNetwork(t *testing.T) {
	type args struct {
		ip      string
		netmask uint8
	}
	
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Standard IP/Network",
			args: args{
				ip: "192.168.1.10",
				netmask: 24,
			},
			want: "192.168.1.0",
			wantErr: false,
		},
		{
			name: "Bad IP Address",
			args: args{
				ip: "541.221.1.10",
				netmask: 24,
			},
			want: "",
			wantErr: true,
		},
		{
			name: "Bad Network Mask",
			args: args{
				ip: "192.168.1.10",
				netmask: 33,
			},
			want: "",
			wantErr: true,
		},
		{
			name: "Bad IP Address (another string)",
			args: args{
				ip: "example.com",
				netmask: 24,
			},
			want: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNetwork(tt.args.ip, tt.args.netmask)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}
