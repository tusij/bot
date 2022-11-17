package utils

import (
	"github.com/tusij/bot.git/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		file   string
		config interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "aaaa", args: struct {
			file   string
			config interface{}
		}{file: "../config.yaml", config: &config.BizConfig{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.file, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getConfigPath(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test_file_path", args: struct{ name string }{name: "config.yaml"}, want: "C:/Users/赖志宇/GolandProjects/bot/utils/config.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigPath(tt.args.name); got != tt.want {
				t.Errorf("getConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
