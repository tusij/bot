package messageScoket

import (
	"github.com/tusij/bot.git/errs"
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/token"
	"log"
	"time"
)

// 职责是为了管理session从而启动监听
type SessionManager interface {
	Start(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error
}

type ChannelSessionManagerImpl struct {
	C chan dto.Session
}

func (c *ChannelSessionManagerImpl) Start(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) error {
	shard := apInfo.Shards
	//保证chanl 至少有一个缓冲
	if shard < 1 {
		shard = 1
	}

	if err := CheckSessionLimit(apInfo); err != nil {
		return err
	}

	startInterval := CalcInterval(apInfo.SessionStartLimit.MaxConcurrency)
	c.C = make(chan dto.Session, shard)
	for i := 0; i < int(shard); i++ {
		session := dto.Session{
			ID:      "",
			URL:     apInfo.URL,
			Token:   *token,
			Intent:  *intents,
			LastSeq: 0,
			Shards: dto.ShardConfig{
				ShardID:    uint32(i),
				ShardCount: shard,
			},
		}
		c.C <- session
	}

	for session := range c.C {
		//控制每两秒执行多少次并发
		time.Sleep(startInterval)
		go c.connect(session)
	}
	return nil
}

func (c *ChannelSessionManagerImpl) connect(session dto.Session) {
	var wsClient WebSocket
	var err error
	wsClient = NewWebSocketClient(session)
	if err = wsClient.Connect(); err != nil {
		log.Printf("ws connect error session:%v err:%v", session, err)
		//如果连接失败需要重试 todo 需要判断URL为空的情况 这种情况不应该重试
		if session.URL != "" {
			c.C <- session
		}
		return
	}

	if session.ID != "" {
		err = wsClient.Resume()
	} else {
		//必须需要鉴权 切告知ws intent 才可以监听到对应的消息
		err = wsClient.Identify()
	}

	if err != nil {
		log.Printf("ChannelSessionManagerImpl.connect resume or indentify error session:%v error:%v", session, err)
		return
	}

	if err = wsClient.Listening(); err != nil {
		e := errs.Error(err)
		s := *wsClient.Session()
		//session 可能因为超时导致无法重连 这时候需要释放掉 session id
		if e.Code() == errs.CodeConnCloseCantResume {
			s.ID = ""
			s.LastSeq = 0
		}
		c.C <- s
	}
}
