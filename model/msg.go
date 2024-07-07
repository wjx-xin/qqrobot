package model

import "time"

// 作者信息
type Author struct {
	Avatar   string `json:"avatar,omitempty"` // 使用omitempty
	Bot      bool   `json:"bot,omitempty"`
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

// 成员信息
type Member struct {
	JoinedAt time.Time `json:"joined_at,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Roles    []string  `json:"roles,omitempty"`
}

// 用户信息
type Mention struct {
	Avatar   string `json:"avatar,omitempty"`
	Bot      bool   `json:"bot,omitempty"`
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

// JSON结构
type MessageData struct {
	OpCode   int    `json:"op,omitempty"`
	Sequence int    `json:"s,omitempty"`
	Type     string `json:"t,omitempty"`
	ID       string `json:"id,omitempty"`
	Data     *Data  `json:"d,omitempty"`
}

// Data 内嵌结构体，使用 "-" 作为别名表示即使字段存在也不解析
type Data struct {
	Author            *Author    `json:"author,omitempty"`
	ChannelID         string     `json:"channel_id,omitempty"`
	Content           string     `json:"content,omitempty"`
	GuildID           string     `json:"guild_id,omitempty"`
	MessageID         string     `json:"id,omitempty"`
	Member            *Member    `json:"member,omitempty"`
	Mentions          []*Mention `json:"mentions,omitempty"`
	Sequence          int        `json:"seq,omitempty"`
	SequenceInChannel string     `json:"seq_in_channel,omitempty"`
	Timestamp         time.Time  `json:"timestamp,omitempty"`
}
