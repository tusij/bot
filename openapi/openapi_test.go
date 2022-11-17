package openapi

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/token"
	"testing"
	"time"
)

func TestOpenAPImpl_GetWSInfo(t *testing.T) {
	type fields struct {
		token       *token.Token
		timeout     time.Duration
		idleConns   int
		lastTraceID string
		restyClient *resty.Client
	}
	type args struct {
		ctx    context.Context
		params map[string]string
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.WebsocketAP
		wantErr bool
	}{
		{
			name: "websocket get ws info test",
			fields: fields{
				token: &token.Token{
					AppID:       102031290,
					AccessToken: "zPzcRITV8ggJmAIYmlQXTc1iN0GTaRNV",
					Type:        "Bot",
				},
				timeout:     3000 * time.Millisecond,
				idleConns:   100,
				lastTraceID: "",
				restyClient: nil,
			},
			args: args{
				ctx:    context.Background(),
				params: nil,
				body:   "",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OpenAPImpl{
				token:       tt.fields.token,
				timeout:     tt.fields.timeout,
				idleConns:   tt.fields.idleConns,
				lastTraceID: tt.fields.lastTraceID,
				restyClient: tt.fields.restyClient,
			}
			o.setupClient()
			got, err := o.GetWSInfo(tt.args.ctx, tt.args.params, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWSInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

func TestOpenAPImpl_PostMessage(t *testing.T) {
	type fields struct {
		token       *token.Token
		timeout     time.Duration
		idleConns   int
		lastTraceID string
		restyClient *resty.Client
	}
	type args struct {
		ctx       context.Context
		channelID string
		msg       *dto.MessageToCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.Message
		wantErr bool
	}{
		{
			name: "websocket post message test",
			fields: fields{
				token: &token.Token{
					AppID:       102031290,
					AccessToken: "zPzcRITV8ggJmAIYmlQXTc1iN0GTaRNV",
					Type:        "Bot",
				},
				timeout:     3000 * time.Millisecond,
				idleConns:   100,
				lastTraceID: "",
				restyClient: nil,
			},
			args: args{
				ctx:       context.Background(),
				channelID: "13677259",
				msg: &dto.MessageToCreate{
					Content: "测试",
					MessageReference: &dto.MessageReference{
						MessageID:             "",
						IgnoreGetMessageError: true,
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OpenAPImpl{
				token:       tt.fields.token,
				timeout:     tt.fields.timeout,
				idleConns:   tt.fields.idleConns,
				lastTraceID: tt.fields.lastTraceID,
				restyClient: tt.fields.restyClient,
			}
			o.setupClient()
			got, err := o.PostMessage(tt.args.ctx, tt.args.channelID, tt.args.msg)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
