// Package app provides web application that can handle HTTP requests

package app

import (
	"os"
	"testing"
	"time"
)

func Test_getIntEnv(t *testing.T) {
	type args struct {
		envName      string
		defaultValue int
	}
	os.Setenv("SomeIntEnv", "1")
	os.Setenv("SomeStrnv", "hello")
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Environment variable",
			args: args{
				envName: "SomeIntEnv",
				defaultValue: 10,
			},
			want: 1,
		},
		{
			name: "Haven't environment variable",
			args: args{
				envName: "SomeTestEnv1",
				defaultValue: 10,
			},
			want: 10,
		},
		{
			name: "Not int environment variable",
			args: args{
				envName: "SomeStrnv",
				defaultValue: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIntEnv(tt.args.envName, tt.args.defaultValue); got != tt.want {
				t.Errorf("getIntEnv() = %v, want %v", got, tt.want)
			}
		})
	}
	os.Unsetenv("SomeIntEnv")
	os.Unsetenv("SomeStrnv")
}

func Test_getVars(t *testing.T) {
	os.Setenv("LIMIT_PER_MINUTE", "100")
	os.Setenv("COOLDOWN_PERIOD_IN_MINUTES", "1")
	os.Setenv("CLEAN_PERIOD_IN_MINUTES", "1")
	os.Setenv("NETMASK", "24")
	tests := []struct {
		name               string
		wantLimitPerMinute int
		wantNetmask        uint8
		wantCooldownPeriod time.Duration
		wantCleanPeriod    time.Duration
	}{
		{
			name: "Default check",
			wantLimitPerMinute: 100,
			wantNetmask: 24,
			wantCooldownPeriod: 1 * time.Minute,
			wantCleanPeriod: 1 * time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimitPerMinute, gotNetmask, gotCooldownPeriod, gotCleanPeriod := getVars()
			if gotLimitPerMinute != tt.wantLimitPerMinute {
				t.Errorf("getVars() gotLimitPerMinute = %v, want %v", gotLimitPerMinute, tt.wantLimitPerMinute)
			}
			if gotNetmask != tt.wantNetmask {
				t.Errorf("getVars() gotNetmask = %v, want %v", gotNetmask, tt.wantNetmask)
			}
			if gotCooldownPeriod != tt.wantCooldownPeriod {
				t.Errorf("getVars() gotCooldownPeriod = %v, want %v", gotCooldownPeriod, tt.wantCooldownPeriod)
			}
			if gotCleanPeriod != tt.wantCleanPeriod {
				t.Errorf("getVars() gotCleanPeriod = %v, want %v", gotCleanPeriod, tt.wantCleanPeriod)
			}
		})
	}
	os.Unsetenv("LIMIT_PER_MINUTE")
	os.Unsetenv("COOLDOWN_PERIOD_IN_MINUTES")
	os.Unsetenv("CLEAN_PERIOD_IN_MINUTES")
	os.Unsetenv("NETMASK")
}
