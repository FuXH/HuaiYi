package chat_task

import (
	"chat_service/tool_function"
	"fmt"
	"time"

	hyentity "chat_service/entity/hunyuan_msg"
	roleentity "chat_service/entity/role"
	"chat_service/repository/remote/hunyuan"
)

type ChatTask struct {
	role           *roleentity.Role       // 角色描述
	latestChatTime time.Time              // 最后一次对话的时间
	chatContext    []*hyentity.HyMessage  // 对话上下文
	hyClient       *hunyuan.HyCli         // 混元client
	llmConfig      *hyentity.HyChatConfig // llm对话配置
	funcCallList   []*hyentity.HyTool     // 支持的func_call
	memory         []string               // 待定
}

func Init(role *roleentity.Role, hyClient *hunyuan.HyCli,
	llmConfig *hyentity.HyChatConfig,
	funcCallList []*hyentity.HyTool) *ChatTask {
	task := &ChatTask{
		role:           role,
		hyClient:       hyClient,
		llmConfig:      llmConfig,
		funcCallList:   funcCallList,
		latestChatTime: time.Now(),
		chatContext:    make([]*hyentity.HyMessage, 0),
	}
	if role != nil {
		task.chatContext = append(task.chatContext, &hyentity.HyMessage{
			Role:    hyentity.HyRoleSystem,
			Content: role.Desc,
		})
	}
	return task
}

func (p *ChatTask) Exec(input string, args ...interface{}) (*hyentity.HyChatRsp, error) {
	reqMsg := &hyentity.HyMessage{
		Role:    hyentity.HyRoleUser,
		Content: input,
	}
	chatRsp, err := p.hyClient.Chat(append(p.chatContext, reqMsg), p.llmConfig, p.funcCallList)
	if err != nil {
		return nil, fmt.Errorf("ChatTask.Exec func(Chat) fail, err: %v", err)
	} else if len(chatRsp.Choices) == 0 {
		return nil, fmt.Errorf("ChatTask.Exec func(Chat) choice is null, rsp: %v", chatRsp)
	}

	choice := chatRsp.Choices[0]
	rspMsg := &hyentity.HyMessage{
		Role:    hyentity.HyRoleAssistant,
		Content: choice.Message.Content,
	}
	p.chatContext = append(p.chatContext, reqMsg, rspMsg)

	// 需要调用func_call
	if p.needFuncCall(chatRsp) {
		toolMsg, err := p.callTool(chatRsp)
		if err != nil {
			return nil, fmt.Errorf("ChatTask.Exec func(callTool) fail, err: %v", err)
		}
		chatRsp, err = p.hyClient.Chat(append(p.chatContext, toolMsg), p.llmConfig, p.funcCallList)
		if err != nil {
			return nil, fmt.Errorf("ChatTask.Exec.needFuncCall func(Chat) fail, err: %v", err)
		}
	}

	return chatRsp, nil
}

func (p *ChatTask) Query(param string) interface{} {
	switch param {
	default:

	}
	return nil
}

func (p *ChatTask) needFuncCall(chatRsp *hyentity.HyChatRsp) bool {
	if chatRsp.Choices[0].FinishReason != hyentity.FinishReasonTool {
		return false
	}
	if len(chatRsp.GetToolCalls(p.llmConfig.IsStream)) == 0 {
		return false
	}
	return true
}

func (p *ChatTask) callTool(chatRsp *hyentity.HyChatRsp) (*hyentity.HyMessage, error) {
	// 调用tool
	tool := chatRsp.GetToolCalls(p.llmConfig.IsStream)[0]
	toolRsp := ""
	for _, val := range p.funcCallList {
		if val.Function.Name == tool.Function.Name {
			toolRsp = tool_function.CallFunction(val.Function.Name, "")
			fmt.Printf("val.Function.Name: %s, toolRsp: %s\n", val.Function.Name, toolRsp)
		}
	}

	// 返回大模型tool信息
	toolMsg := &hyentity.HyMessage{
		Role:       hyentity.HyRoleTool,
		Content:    toolRsp,
		ToolCallId: tool.ID,
		ToolCalls:  chatRsp.GetToolCalls(p.llmConfig.IsStream),
	}
	return toolMsg, nil
}
