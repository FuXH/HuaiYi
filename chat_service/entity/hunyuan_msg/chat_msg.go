package hunyuan_msg

import "time"

// HyChatConfig 对话的配置
type HyChatConfig struct {
	IsStream    bool
	RspChan     chan string // 返回chan
	Temperature float64     // 取值区间为 [0.0, 2.0]
	ToolChoice  string      // none/auto/custom

	// 非流式

	// 流式
	GapTime    time.Duration // 流式返回的间隔
	RetryCount int           // 获取空数据的等待次数
}

func NewChatConfig() *HyChatConfig {
	return &HyChatConfig{
		IsStream:   false,
		RspChan:    make(chan string),
		ToolChoice: "auto",
		GapTime:    time.Second,
		RetryCount: 1,
	}
}
