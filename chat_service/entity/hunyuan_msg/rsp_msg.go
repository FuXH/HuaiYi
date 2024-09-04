package hunyuan_msg

import "fmt"

const (
	FinishReasonStop = "stop"
	FinishReasonTool = "tool_calls"
)

// HyChatRsp 对话的返回体
type HyChatRsp struct {
	ID        string          `json:"Id"`
	Note      string          `json:"Note"`
	Choices   []*HyChatChoice `json:"Choices"`
	RequestID string          `json:"RequestId"`
}

type HyChatChoice struct {
	Message      HyChatMessage `json:"Message"`
	Delta        HyChatDelta   `json:"Delta"`
	FinishReason string        `json:"FinishReason"`
}

type HyChatMessage struct {
	Role      string        `json:"Role"`
	Content   string        `json:"Content"`
	ToolCalls []*HyToolCall `json:"ToolCalls"`
}

type HyChatDelta struct {
	Role      HyRole        `json:"Role"`
	Content   string        `json:"Content"`
	ToolCalls []*HyToolCall `json:"ToolCalls"`
}

func (p *HyChatRsp) GetContent(isStream bool) string {
	if len(p.Choices) == 0 {
		return ""
	}

	choice := p.Choices[0]
	content := ""
	if isStream {

	} else {
		content = choice.Message.Content
	}
	return content
}

func (p *HyChatRsp) GetToolCalls(isStream bool) []*HyToolCall {
	if len(p.Choices) == 0 {
		return nil
	}

	toolCalls := make([]*HyToolCall, 0)
	choice := p.Choices[0]
	if isStream {
		toolCalls = choice.Delta.ToolCalls
	} else {
		toolCalls = choice.Message.ToolCalls
	}
	return toolCalls
}

func (p *HyChatRsp) GetRspContent(isStream bool, rspChan chan (string)) error {
	if len(p.Choices) == 0 {
		return fmt.Errorf("GetRspContent is empty")
	}

	choice := p.Choices[0]
	if isStream {

	} else {
		rspChan <- choice.Message.Content
	}
	return nil
}

func (p *HyChatRsp) Display() {
	fmt.Println("ID: ", p.ID)
	fmt.Println("Note: ", p.Note)
	for _, val := range p.Choices {
		fmt.Println("  FinishReason: ", val.FinishReason)
		fmt.Println("    Choice Message: ", val.Message.Role)
		fmt.Println("    Choice Message: ", val.Message.Content)
		for _, tool := range val.Message.ToolCalls {
			fmt.Println("    Choice Message Tool: ", tool.ID)
			fmt.Println("    Choice Message Tool: ", tool.Type)
			fmt.Println("    Choice Message Tool: ", tool.Function)
		}
		fmt.Println("    Choice Delta: ", val.Message.Role)
		fmt.Println("    Choice Delta: ", val.Message.Content)
		for _, tool := range val.Message.ToolCalls {
			fmt.Println("    Choice Delta Tool: ", tool.ID)
			fmt.Println("    Choice Delta Tool: ", tool.Type)
			fmt.Println("    Choice Delta Tool: ", tool.Function)
		}
	}
}
