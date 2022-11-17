package messageScoket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tusij/bot.git/errs"
	"github.com/tusij/bot.git/event"
	"github.com/tusij/bot.git/modle/dto"
	"log"
	"time"
)

const (
	TEXT   = 1
	BINARY = 2
	CLOSE  = 8
	PING   = 9
	PONG   = 10
)

type WebSocket interface {
	// Connect 连接到 wss 地址
	Connect() error
	// 鉴权连接 连接鉴权后才能使用
	Identify() error

	Session() *dto.Session
	//重连
	Resume() error
	// 监听websocket事件
	Listening() error
	//发送数据 主要为心跳 重连等
	Write(message *dto.WSPayload) error
	// Close 关闭连接
	Close()
}

type WebSocketDefaultImpl struct {
	conn            *websocket.Conn
	messageQueue    chan *dto.WSPayload
	session         *dto.Session
	errChan         chan error
	heartBeatTicker *time.Ticker
}

func NewWebSocketClient(session dto.Session) WebSocket {
	return &WebSocketDefaultImpl{
		messageQueue:    make(chan *dto.WSPayload, 100),
		session:         &session,
		errChan:         make(chan error, 10),
		heartBeatTicker: time.NewTicker(30 * time.Second),
	}
}

func (w *WebSocketDefaultImpl) Connect() error {
	if w.session.URL == "" {
		return errs.ErrURLInvalid
	}
	conn, _, err := websocket.DefaultDialer.Dial(w.session.URL, nil)
	if err != nil {
		return err
	}
	w.conn = conn
	return nil
}

func (w *WebSocketDefaultImpl) Identify() error {
	payload := &dto.WSPayload{
		Data: &dto.WSIdentityData{
			Token:   w.session.Token.GetString(),
			Intents: w.session.Intent,
			Shard: []uint32{
				w.session.Shards.ShardID,
				w.session.Shards.ShardCount,
			},
		},
	}
	payload.OPCode = dto.WSIdentity
	return w.Write(payload)
}
func (w *WebSocketDefaultImpl) Session() *dto.Session {
	return w.session
}
func (w *WebSocketDefaultImpl) Resume() error {
	payload := &dto.WSPayload{
		Data: &dto.WSResumeData{
			Token:     w.session.Token.GetString(),
			SessionID: w.session.ID,
			Seq:       w.session.LastSeq,
		},
	}
	payload.OPCode = dto.WSResume // 内嵌结构体字段，单独赋值
	return w.Write(payload)
}

// 监听websocket事件
func (w *WebSocketDefaultImpl) Listening() error {
	go w.acceptMessage()
	go w.handleMessage()
	// handler message
	for {
		select {
		case err := <-w.errChan:

			log.Printf("%s Listening stop. err is %v \n", w.session, err)
			//这个错误指机器人无权限 session不会再被重连
			if errs.Error(err).Code() == 4196 || errs.Error(err).Code() == 4915 {
				err = errs.New(errs.CodeConnCloseCantIdentify, err.Error())
			}
			return err
		//	心跳 keepalive
		case <-w.heartBeatTicker.C:
			heartBeatEvent := &dto.WSPayload{
				WSPayloadBase: dto.WSPayloadBase{
					OPCode: dto.WSHeartbeat,
				},
				Data: w.session.LastSeq,
			}
			_ = w.Write(heartBeatEvent)
		}
	}
	return nil
}

func (w *WebSocketDefaultImpl) acceptMessage() {
	defer close(w.messageQueue)

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			w.errChan <- err
			return
		}

		payload := &dto.WSPayload{}
		if err = json.Unmarshal(message, payload); err != nil {
			log.Printf("WebSocketDefaultImpl.acceptMessage json unmarshal error message:%s err:%v", message, err)
			continue
		}
		payload.RawMessage = message
		if w.isHandleBuildIn(payload) {
			continue
		}
		//fmt.Println(payload)
		w.messageQueue <- payload
	}
}

func (w *WebSocketDefaultImpl) handleMessage() {
	defer func() {
		//意外情况下 将error抛出 由上游决定是否需要重连或者退出
		if err := recover(); err != nil {
			w.errChan <- fmt.Errorf("panic: %v", err)
		}
	}()
	for payload := range w.messageQueue {
		if payload.Seq > 0 {
			w.session.LastSeq = payload.Seq
		}

		// 鉴权后会ws会接收到ready信息 这时候需要补充 session id等信息回client中
		if payload.Type == "READY" {
			w.readyHandler(payload)
			continue
		}
		// 解析具体事件，并投递给业务注册的 handler
		if err := event.ParseAndHandle(payload); err != nil {
			log.Printf("%v parseAndHandle failed, %v", payload, err)
		}
	}
}

func (w *WebSocketDefaultImpl) readyHandler(payload *dto.WSPayload) {
	readyData := &dto.WSReadyData{}
	if err := event.ParseData(payload.RawMessage, readyData); err != nil {
		log.Printf("parseReadyData failed, session:%v, err:%v, message %v \n", w.session, err, payload.RawMessage)
	}

	// 基于 ready 事件，更新 session 信息
	w.session.ID = readyData.SessionID
	w.session.Shards.ShardID = readyData.Shard[0]
	w.session.Shards.ShardCount = readyData.Shard[1]

}

func (w *WebSocketDefaultImpl) isHandleBuildIn(payload *dto.WSPayload) bool {
	switch payload.OPCode {
	case dto.WSHello: // 接收到 hello 后需要开始发心跳
		func(message []byte) {
			helloData := &dto.WSHelloData{}
			if err := event.ParseData(message, helloData); err != nil {
				log.Printf("WebSocketDefaultImpl.isHandleBuildIn heart beat error session:%v message:%v err:%v\n",
					w.session, message, err)
			}
			// 根据 hello 的回包，重新设置心跳的定时器时间
			w.heartBeatTicker.Reset(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
		}(payload.RawMessage)
	case dto.WSHeartbeatAck: // 心跳 ack 不需要业务处理
	case dto.WSReconnect: // 达到连接时长，需要重新连接，此时可以通过 resume 续传原连接上的事件
		w.errChan <- errs.ErrNeedReConnect
	case dto.WSInvalidSession: // 无效的 sessionLog，需要重新鉴权
		w.errChan <- errs.ErrInvalidSession
	default:
		return false
	}
	return true
}

//发送数据 主要为心跳 重连等
func (w *WebSocketDefaultImpl) Write(message *dto.WSPayload) error {
	m, err := json.Marshal(message)
	if err != nil {
		log.Printf("WebSocketDefaultImpl.Write json marshal error messgae:%v err:%v", message, err)
	}

	if err = w.conn.WriteMessage(TEXT, m); err != nil {
		log.Printf("WebSocketDefaultImpl.Write write message error messgae:%v err:%v", message, err)
		w.errChan <- err
		return err
	}
	return nil
}
func (w *WebSocketDefaultImpl) Close() {
	err := w.conn.Close()
	if err != nil {
		log.Printf("conn close fail, session:%v, err:%v\n", w.session, err)
	}
	w.heartBeatTicker.Stop()
}
