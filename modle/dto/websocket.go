package dto

import (
	"fmt"
	"github.com/tusij/bot.git/token"
)

// WebsocketAP wss 接入点信息
type WebsocketAP struct {
	URL               string            `json:"url"`
	Shards            uint32            `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

// SessionStartLimit 链接频控信息
type SessionStartLimit struct {
	Total          uint32 `json:"total"`
	Remaining      uint32 `json:"remaining"`
	ResetAfter     uint32 `json:"reset_after"`
	MaxConcurrency uint32 `json:"max_concurrency"`
}

// ShardConfig 连接的 shard 配置，ShardID 从 0 开始，ShardCount 最小为 1
type ShardConfig struct {
	ShardID    uint32
	ShardCount uint32
}

// Session 连接的 session 结构，包括链接的所有必要字段
type Session struct {
	ID      string
	URL     string
	Token   token.Token
	Intent  Intent
	LastSeq uint32
	Shards  ShardConfig
}

// String 输出session字符串
func (s *Session) String() string {
	return fmt.Sprintf("[ws][ID:%s][Shard:(%d/%d)]",
		s.ID, s.Shards.ShardID, s.Shards.ShardCount, s.Intent)
}

type Intent int

// websocket intent 声明
const (
	// IntentGuilds 包含
	// - GUILD_CREATE
	// - GUILD_UPDATE
	// - GUILD_DELETE
	// - GUILD_ROLE_CREATE
	// - GUILD_ROLE_UPDATE
	// - GUILD_ROLE_DELETE
	// - CHANNEL_CREATE
	// - CHANNEL_UPDATE
	// - CHANNEL_DELETE
	// - CHANNEL_PINS_UPDATE
	IntentGuilds Intent = 1 << iota

	// IntentGuildMembers 包含
	// - GUILD_MEMBER_ADD
	// - GUILD_MEMBER_UPDATE
	// - GUILD_MEMBER_REMOVE
	IntentGuildMembers

	IntentGuildBans
	IntentGuildEmojis
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages

	// IntentGuildMessageReactions 包含
	// - MESSAGE_REACTION_ADD
	// - MESSAGE_REACTION_REMOVE
	IntentGuildMessageReactions

	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping

	IntentInteraction Intent = 1 << 26 // 互动事件
	IntentAudit       Intent = 1 << 27 // 审核事件
	// IntentForum 论坛事件
	//  - THREAD_CREATE     // 当用户创建主题时
	//  - THREAD_UPDATE     // 当用户更新主题时
	//  - THREAD_DELETE     // 当用户删除主题时
	//  - POST_CREATE       // 当用户创建帖子时
	//  - POST_DELETE       // 当用户删除帖子时
	//  - REPLY_CREATE      // 当用户回复评论时
	//  - REPLY_DELETE      // 当用户回复评论时
	//  - FORUM_PUBLISH_AUDIT_RESULT      // 当用户发表审核通过时
	IntentForum Intent = 1 << 28 // 论坛事件

	// IntentAudio
	//  - AUDIO_START           // 音频开始播放时
	//  - AUDIO_FINISH          // 音频播放结束时
	IntentAudio          Intent = 1 << 29 // 音频机器人事件
	IntentGuildAtMessage Intent = 1 << 30 // 只接收@消息事件

	IntentNone Intent = 0
)
