package hunyuan_msg

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

type HyRole string

const (
	HyRoleSystem    HyRole = "system"
	HyRoleUser      HyRole = "user"
	HyRoleAssistant HyRole = "assistant"
	HyRoleTool      HyRole = "tool"
)

// HyMessage 回复的结构体
type HyMessage struct {
	Role       HyRole        // 角色
	Content    string        // 文本内容
	ToolCallId string        // tool调用结果，当role=tool时使用
	ToolCalls  []*HyToolCall // 调用tool
}

// HyToolCall 函数调用说明
type HyToolCall struct {
	ID       string
	Type     string // 目前只支持function
	Function *HyToolCallFunction
}

type HyToolCallFunction struct {
	Name      string // function名称
	Arguments string // function参数
}

func (s HyToolCall) Convert() *hunyuan.ToolCall {
	res := &hunyuan.ToolCall{
		Id:   common.StringPtr(s.ID),
		Type: common.StringPtr(s.Type),
		Function: &hunyuan.ToolCallFunction{
			Name:      common.StringPtr(s.Function.Name),
			Arguments: common.StringPtr(s.Function.Arguments),
		},
	}
	return res
}
