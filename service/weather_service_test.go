package service

import (
	"context"
	"fmt"
	"testing"
)

func TestWeatherServiceJuHeImpl_QueryCityWeather(t *testing.T) {
	s := NewJuHeWeatherService(3000, 100)
	ctx := context.Background()
	weather, err := s.QueryCityWeather(ctx, "深圳")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(weather)
}
