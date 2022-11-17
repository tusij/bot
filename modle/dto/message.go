package dto

// Message 消息结构体定义
type Message struct {
	// 消息ID
	ID string `json:"id"`
	// 子频道ID
	ChannelID string `json:"channel_id"`
	// 频道ID
	GuildID string `json:"guild_id"`
	// 内容
	Content string `json:"content"`
	// 发送时间
	Timestamp string `json:"timestamp"`
	// 消息编辑时间
	EditedTimestamp string `json:"edited_timestamp"`
	// 是否@all
	MentionEveryone bool `json:"mention_everyone"`
	// 消息发送方

	// 附件
	Attachments []*MessageAttachment `json:"attachments"`
	// 结构化消息-embeds
	Embeds []*Embed `json:"embeds"`

	// 私信消息
	DirectMessage bool `json:"direct_message"`
	// 子频道 seq，用于消息间的排序，seq 在同一子频道中按从先到后的顺序递增，不同的子频道之前消息无法排序
	SeqInChannel string `json:"seq_in_channel"`

	// 私信场景下，该字段用来标识从哪个频道发起的私信
	SrcGuildID string `json:"src_guild_id"`
}

// Embed 结构
type Embed struct {
	Title       string                `json:"title,omitempty"`
	Description string                `json:"description,omitempty"`
	Prompt      string                `json:"prompt"` // 消息弹窗内容，消息列表摘要
	Thumbnail   MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Fields      []*EmbedField         `json:"fields,omitempty"`
}

// MessageEmbedThumbnail embed 消息的缩略图对象
type MessageEmbedThumbnail struct {
	URL string `json:"url"`
}

// EmbedField Embed字段描述
type EmbedField struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// MessageAttachment 附件定义
type MessageAttachment struct {
	URL string `json:"url"`
}
