package service

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/tusij/bot.git/modle/dto"
	"net"
	"net/http"
	"time"
)

const (
	token       = "3d9b19211b44ce2b8f4d9cbd9ca4d582"
	weather_api = "http://apis.juhe.cn/simpleWeather/query"
)

type WeatherService interface {
	QueryCityWeather(ctx context.Context, city string) (*dto.Weather, error)
}

type WeatherServiceJuHeImpl struct {
	token     string
	timeout   time.Duration
	idleConns int

	restyClient *resty.Client // resty client http客户端工具
}

func NewJuHeWeatherService(timeoutMS, idleConns int) WeatherService {
	s := &WeatherServiceJuHeImpl{
		token:     token,
		timeout:   time.Duration(timeoutMS) * time.Millisecond,
		idleConns: idleConns,

		restyClient: nil,
	}

	s.setupClient()
	return s
}

func createTransport(localAddr net.Addr, idleConns int) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	if localAddr != nil {
		dialer.LocalAddr = localAddr
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          idleConns,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   idleConns,
		MaxConnsPerHost:       idleConns,
	}
}

// 初始化 client
func (w *WeatherServiceJuHeImpl) setupClient() {
	w.restyClient = resty.New().
		SetTransport(createTransport(nil, w.idleConns)). // 自定义 transport
		SetTimeout(w.timeout)
}

func (w *WeatherServiceJuHeImpl) QueryCityWeather(ctx context.Context, city string) (*dto.Weather, error) {
	resp, err := w.restyClient.R().SetContext(ctx).
		SetResult(dto.WeatherResponse{}).
		SetQueryParams(map[string]string{
			"key":  w.token,
			"city": city,
		}).
		Get(weather_api)
	if err != nil {
		return nil, err
	}
	return &resp.Result().(*dto.WeatherResponse).Result.Realtime, nil
}
