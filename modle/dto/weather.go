package dto

import "fmt"

type Weather struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	Info        string `json:"info"`
	Wid         string `json:"wid"`
	Direct      string `json:"direct"`
	Power       string `json:"power"`
	Aqi         string `json:"aqi"`
}

func (w *Weather) ToString() string {
	return fmt.Sprintf("温度:%s\n湿度:%s\n天气情况:%s\n风向:%s\n风力:%s\n空气质量:%s\n",
		w.Temperature, w.Humidity, w.Info, w.Direct, w.Power, w.Aqi)
}

type WeatherResponse struct {
	Reason string `json:"reason"`
	Result struct {
		City     string  `json:"city"`
		Realtime Weather `json:"realtime"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}
