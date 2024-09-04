package hunyuan_msg

import (
	"chat_service/tool_function"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

func NewHyTool(funcInfo tool_function.Function) *HyTool {
	fcName, fcDesc, argsDesc := funcInfo.GetInfo()
	return &HyTool{
		Type: "function",
		Function: &HyFunction{
			Name:        fcName,
			Parameters:  argsDesc,
			Description: fcDesc,
		},
	}
}

// HyTool tool_function
type HyTool struct {
	Type     string // 工具类型，当前只支持function
	Function *HyFunction
}

type HyFunction struct {
	Name        string // function名称，只能包含a-z，A-Z，0-9，\_或-
	Parameters  string // 传参描述
	Description string // 函数描述
}

func (s *HyTool) Convert() *hunyuan.Tool {
	res := &hunyuan.Tool{
		Type: common.StringPtr(s.Type),
		Function: &hunyuan.ToolFunction{
			Name:        common.StringPtr(s.Function.Name),
			Parameters:  common.StringPtr(s.Function.Parameters),
			Description: common.StringPtr(s.Function.Description),
		},
	}
	return res
}

func (s *HyTool) Call() string {
	toolFunc := tool_function.FunctionList[s.Function.Name]
	if toolFunc == nil {
		return ""
	}
	return toolFunc.Call(s.Function.Parameters)
}
