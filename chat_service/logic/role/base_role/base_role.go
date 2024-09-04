package base_role

import (
	hyentity "chat_service/entity/hunyuan_msg"
	"fmt"
	"time"

	"chat_service/repository/remote/hunyuan"
)

type BaseRole struct {
	RoleName      string // 角色名称
	PromptTplFile string // 模板文件
	RoleDesc      string // 角色描述
	HyClient      *hunyuan.HyCli
	ToolFunction  map[string]*hyentity.HyTool      // 可调用能力，key: funcName
	ChatContext   map[string][]*hyentity.HyMessage // 对话记录，key: chatID
}

func (p *BaseRole) Chat(chatID string,
	msg *hyentity.HyMessage, chatCfg *hyentity.HyChatConfig) (*hyentity.HyChatRsp, error) {
	chatContext := p.ChatContext[chatID]
	if chatID == "" {
		chatID = p.GenerateID()
		chatContext = append([]*hyentity.HyMessage{}, &hyentity.HyMessage{
			Role:    hyentity.HyRoleSystem,
			Content: p.RoleDesc,
		})
		p.ChatContext[chatID] = chatContext
	}
	chatContext = append(p.ChatContext[chatID], msg)
	p.ChatContext[chatID] = chatContext

	toolList := make([]*hyentity.HyTool, 0)
	for _, tool := range p.ToolFunction {
		toolList = append(toolList, tool)
	}

	hyRsp, err := p.HyClient.Chat(chatContext, chatCfg, toolList)
	if err != nil {
		return nil, fmt.Errorf("HyClient.Chat fail, err: %v", err)
	}
	if len(hyRsp.Choices) == 0 {
		return nil, fmt.Errorf("HyClient.Chat choice is null, rsp: %v", hyRsp)
	}
	choice := hyRsp.Choices[0]
	p.ChatContext[chatID] = append(p.ChatContext[chatID], &hyentity.HyMessage{
		Role:    hyentity.HyRoleAssistant,
		Content: choice.Message.Content,
	})

	return p.CallTool(chatID, hyRsp, chatCfg)
}

func (p *BaseRole) CallTool(chatID string, hyRsp *hyentity.HyChatRsp, chatCfg *hyentity.HyChatConfig) (*hyentity.HyChatRsp, error) {
	if hyRsp.Choices[0].FinishReason != hyentity.FinishReasonTool {
		return hyRsp, nil
	}
	if len(hyRsp.GetToolCalls(chatCfg.IsStream)) == 0 {
		return hyRsp, nil
	}

	// 调用tool
	tool := hyRsp.GetToolCalls(chatCfg.IsStream)[0]
	toolRsp := p.ToolFunction[tool.Function.Name].Call()

	// 返回大模型tool信息
	toolMsg := &hyentity.HyMessage{
		Role:       hyentity.HyRoleTool,
		Content:    toolRsp,
		ToolCallId: tool.ID,
		ToolCalls:  hyRsp.GetToolCalls(chatCfg.IsStream),
	}
	return p.Chat(chatID, toolMsg, chatCfg)
}

func (p *BaseRole) EditMsg(msg *hyentity.HyMessage) string {
	return msg.Content
}

func (p *BaseRole) Output(chatID string, hyRsp *hyentity.HyChatRsp, chatCfg *hyentity.HyChatConfig) error {
	hyRsp.Display()
	return nil
}

func (p *BaseRole) GenerateID() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
