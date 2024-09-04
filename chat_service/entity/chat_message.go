package entity

import "time"

// Chat 一次会话
type Chat struct {
	ChatContext []*HyMessage
	ChatConfig  *HyChatConfig
	Timeout     time.Duration
}
