package config

import "github.com/tusij/bot.git/token"

type BizConfig struct {
	OpenAPIConfig        *OpenAPIConfig        `yaml:"open_api_config"`
	WeatherServiceConfig *WeatherServiceConfig `yaml:"weather_service_config"`
}

type HttpClientConfig struct {
	Timeout   int `yaml:"timeout"`
	IdleConns int `yaml:"idle_conns"`
}

type WeatherServiceConfig struct {
	Token            string            `yaml:"token"`
	HttpClientConfig *HttpClientConfig `yaml:"http_client_config"`
}

type OpenAPIConfig struct {
	Token            *token.Token      `yaml:"token"`
	HttpClientConfig *HttpClientConfig `yaml:"http_client_config"`
}
