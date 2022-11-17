package main

import (
	"context"
	"github.com/tusij/bot.git/event"
	"github.com/tusij/bot.git/messageScoket"
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/openapi"
	"github.com/tusij/bot.git/service"
	"github.com/tusij/bot.git/token"
	"log"
	"regexp"
	"strings"
)

var openAPI openapi.OpenAPI
var weatherService service.WeatherService

func main() {
	t := &token.Token{
		AppID:       102031290,
		AccessToken: "zPzcRITV8ggJmAIYmlQXTc1iN0GTaRNV",
		Type:        token.TypeBot,
	}

	openAPI = &openapi.OpenAPImpl{}
	openAPI = openAPI.Setup(t, 3000, 100)
	ctx := context.Background()
	wsInfo, err := openAPI.GetWSInfo(ctx, nil, "")
	if err != nil {
		log.Fatalf("Get WS Info Fail, error:%v", err)
	}

	weatherService = service.NewJuHeWeatherService(3000, 100)
	event.RegisterHandlers(atMessageHandle(weatherService))
	var manager messageScoket.SessionManager
	manager = &messageScoket.ChannelSessionManagerImpl{}
	//intent表示监听的事件 正常来说可以监听多个事件用二进制表示 目前我们只关注@机器人的事件 二进制第30位位1
	var intent = dto.IntentGuildAtMessage
	if err = manager.Start(wsInfo, t, &intent); err != nil {
		log.Fatalf("manager.Start error, error:%v", err)
	}

}

func atMessageHandle(weatherService service.WeatherService) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.Message) error {
		ctx := context.Background()
		data.Content = string(regexp.MustCompile(`<@!\d+>`).ReplaceAll([]byte(data.Content), []byte("")))
		data.Content = strings.Trim(data.Content, " \u00A0")
		contents := strings.Split(data.Content, " ")
		cmd := contents[0]
		//上面对内容的处理可以抽象为一个方法
		var message *dto.MessageToCreate = &dto.MessageToCreate{
			Content: "我是享享",
			MessageReference: &dto.MessageReference{
				// 引用这条消息
				MessageID:             data.ID,
				IgnoreGetMessageError: true,
			},
		}
		switch cmd {
		case "天气":
			//这里可以抽象一个方法
			if len(contents) > 1 {
				weather, err := weatherService.QueryCityWeather(ctx, contents[1])
				if err != nil {
					log.Printf("Query Weather Error %v\n", err)
					break
				}
				message.Content = weather.ToString()
			}
		default:

		}
		_, err := openAPI.PostMessage(ctx, data.ChannelID, message)
		if err != nil {
			return err
		}
		return nil
	}
}
