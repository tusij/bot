package openapi

import (
	"context"
	"fmt"
	"github.com/tusij/bot.git/config"
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/token"
	"github.com/tusij/bot.git/utils"

	"time"

	"github.com/go-resty/resty/v2"
)

const MaxIdleConns = 3000

type OpenAPImpl struct {
	token       *token.Token
	timeout     time.Duration
	idleConns   int
	lastTraceID string // lastTraceID id 主要用于日志查询

	restyClient *resty.Client // resty client http客户端工具
}

// TraceID 获取 lastTraceID id
func (o *OpenAPImpl) TraceID() string {
	return o.lastTraceID
}

// Setup 生成一个实例
func (o *OpenAPImpl) Setup(c *config.OpenAPIConfig) OpenAPI {
	idleConns := c.HttpClientConfig.IdleConns
	if idleConns == 0 {
		idleConns = MaxIdleConns
	}

	api := &OpenAPImpl{
		token:     c.Token,
		timeout:   time.Duration(c.HttpClientConfig.Timeout) * time.Millisecond,
		idleConns: idleConns,
	}

	api.setupClient() // 初始化可复用的 client
	return api
}

// WithTimeout 设置请求接口超时时间
func (o *OpenAPImpl) WithTimeout(duration time.Duration) OpenAPI {
	o.restyClient.SetTimeout(duration)
	return o
}

// Transport 透传请求
func (o *OpenAPImpl) Transport(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	resp, err := o.request(ctx).SetBody(body).Execute(method, url)
	return resp.Body(), err
}

// 初始化 client
func (o *OpenAPImpl) setupClient() {
	o.restyClient = resty.New().
		SetTransport(utils.CreateTransport(nil, o.idleConns)). // 自定义 transport
		SetTimeout(o.timeout).
		SetAuthToken(o.token.GetString()).
		SetAuthScheme(o.token.Type)
}

// request 每个请求，都需要创建一个 request
func (o *OpenAPImpl) request(ctx context.Context) *resty.Request {
	return o.restyClient.R().SetContext(ctx)
}

func (o *OpenAPImpl) GetWSInfo(ctx context.Context, params map[string]string, body string) (*dto.WebsocketAP, error) {
	response, err := o.request(ctx).
		SetResult(dto.WebsocketAP{}).
		Get(getURL(gatewayBotURI))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return response.Result().(*dto.WebsocketAP), nil

}

func (o *OpenAPImpl) PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate) (*dto.Message, error) {
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg).
		Post(getURL(messagesURI))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}
