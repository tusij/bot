package service

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/tusij/bot.git/modle/dto"
	"testing"
	"time"
)

func TestWeatherServiceJuHeImpl_QueryCityWeather(t *testing.T) {
	type fields struct {
		token       string
		timeout     time.Duration
		idleConns   int
		restyClient *resty.Client
	}
	type args struct {
		ctx  context.Context
		city string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.Weather
		wantErr bool
	}{
		{"test_weather", fields{
			token:       "3d9b19211b44ce2b8f4d9cbd9ca4d582",
			timeout:     3000 * time.Millisecond,
			idleConns:   100,
			restyClient: nil,
		}, args{
			ctx:  context.Background(),
			city: "深圳",
		}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WeatherServiceJuHeImpl{
				token:       tt.fields.token,
				timeout:     tt.fields.timeout,
				idleConns:   tt.fields.idleConns,
				restyClient: tt.fields.restyClient,
			}
			w.setupClient()
			got, err := w.QueryCityWeather(tt.args.ctx, tt.args.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryCityWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
