package messageScoket

import (
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/token"
	"testing"
)

func TestWebSocketDefaultImpl_Connect(t *testing.T) {
	type fields struct {
		session *dto.Session
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestWebSocketDefaultImpl_Connect",
			fields: fields{session: &dto.Session{
				ID:  "",
				URL: "wss://api.sgroup.qq.com/websocket",
				Token: token.Token{
					AppID:       102031290,
					AccessToken: "zPzcRITV8ggJmAIYmlQXTc1iN0GTaRNV",
					Type:        "Bot",
				},
				Intent:  dto.IntentGuildAtMessage,
				LastSeq: 0,
				Shards:  dto.ShardConfig{},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWebSocketClient(*tt.fields.session)
			if err := w.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebSocketDefaultImpl_Identify(t *testing.T) {
	type fields struct {
		session *dto.Session
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestWebSocketDefaultImpl_Connect",
			fields: fields{session: &dto.Session{
				ID:  "",
				URL: "wss://api.sgroup.qq.com/websocket",
				Token: token.Token{
					AppID:       102031290,
					AccessToken: "zPzcRITV8ggJmAIYmlQXTc1iN0GTaRNV",
					Type:        "Bot",
				},
				Intent:  dto.IntentGuildAtMessage,
				LastSeq: 0,
				Shards:  dto.ShardConfig{},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWebSocketClient(*tt.fields.session)
			err := w.Connect()
			if err != nil {
				t.Errorf("Identify() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := w.Identify(); (err != nil) != tt.wantErr {
				t.Errorf("Identify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebSocketDefaultImpl_Resume(t *testing.T) {
	type fields struct {
		session *dto.Session
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestWebSocketDefaultImpl_Connect",
			fields: fields{session: &dto.Session{
				ID:  "",
				URL: "wss://api.sgroup.qq.com/websocket",
				Token: token.Token{
					AppID:       102031290,
					AccessToken: "zPzcRITV8ggJmAIYmlQXTc1iN0GTaRNV",
					Type:        "Bot",
				},
				Intent:  dto.IntentGuildAtMessage,
				LastSeq: 0,
				Shards:  dto.ShardConfig{},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWebSocketClient(*tt.fields.session)
			w.Connect()
			if err := w.Resume(); (err != nil) != tt.wantErr {
				t.Errorf("Resume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
