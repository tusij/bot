package service

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/tusij/bot.git/config"
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/utils"
	"time"
)

const (
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

func NewJuHeWeatherService(c *config.WeatherServiceConfig) WeatherService {
	idleConns := c.HttpClientConfig.IdleConns
	if idleConns < 0 {
		idleConns = 3000
	}
	s := &WeatherServiceJuHeImpl{
		token:     c.Token,
		timeout:   time.Duration(c.HttpClientConfig.Timeout) * time.Millisecond,
		idleConns: idleConns,

		restyClient: nil,
	}

	s.setupClient()
	return s
}

// 初始化 client
func (w *WeatherServiceJuHeImpl) setupClient() {
	w.restyClient = resty.New().
		SetTransport(utils.CreateTransport(nil, w.idleConns)). // 自定义 transport
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
