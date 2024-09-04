package hunyuan

import (
	"chat_service/util"
	"encoding/json"
	"fmt"
	"time"

	hyentity "chat_service/entity/hunyuan_msg"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
)

const (
	HUNYUAN_MODEL_LITE = "hunyuan_msg-lite"
	HUNYUAN_MODEL_PRO  = "hunyuan-pro"
	HUNYUAN_MODEL_FUNC = "hunyuan-functioncall"
)

var (
	hyCli *HyCli
)

func GetInstance() *HyCli {
	if hyCli == nil {
		secretID := ""
		secretKey := ""
		hyCli = NewHyClient(secretID, secretKey)
	}
	return hyCli
}

// NewHyClient 创建混元请求handler
func NewHyClient(secretID, secretKey string) *HyCli {
	return &HyCli{
		model:     HUNYUAN_MODEL_PRO,
		secretID:  secretID,
		secretKey: secretKey,
	}
}

type HyCli struct {
	model     string
	secretID  string
	secretKey string
}

func (p *HyCli) Chat(msg []*hyentity.HyMessage,
	chatCfg *hyentity.HyChatConfig,
	tools []*hyentity.HyTool) (*hyentity.HyChatRsp, error) {
	// req
	request := hunyuan.NewChatCompletionsRequest()
	request.ToolChoice = common.StringPtr(chatCfg.ToolChoice)
	request.Stream = common.BoolPtr(chatCfg.IsStream)

	for _, val := range msg {
		hyMsg := &hunyuan.Message{
			Role:       common.StringPtr(string(val.Role)),
			Content:    common.StringPtr(val.Content),
			ToolCallId: common.StringPtr(val.ToolCallId),
			ToolCalls:  make([]*hunyuan.ToolCall, 0),
		}
		request.Messages = append(request.Messages, hyMsg)
	}
	for _, tool := range tools {
		request.Tools = append(request.Tools, tool.Convert())
	}
	response, err := p.clientChat(request)
	if err != nil {
		return nil, err
	}

	// rsp
	rspBody := &struct {
		Response *hyentity.HyChatRsp
	}{}
	if chatCfg.IsStream {
		if err := p.chatByStream(request, chatCfg.RspChan, chatCfg.GapTime, chatCfg.RetryCount); err != nil {
			return nil, err
		}

	} else {
		if err = json.Unmarshal([]byte(response.ToJsonString()), rspBody); err != nil {
			return nil, fmt.Errorf("chatByOnce-Unmarshal fail, rspBody: %s, err: %s", response.ToJsonString(), err)
		}
	}

	return rspBody.Response, nil
}

// HyByOnce 非流式调用混元接口
func (p *HyCli) chatByOnce(request *hunyuan.ChatCompletionsRequest) (*hyentity.HyChatRsp, error) {
	response, err := p.clientChat(request)
	if err != nil {
		return nil, fmt.Errorf("chatByOnce-clientChat fail, err: %s", err)
	}
	rspBody := &struct {
		Response *hyentity.HyChatRsp
	}{
		Response: &hyentity.HyChatRsp{},
	}
	if err = json.Unmarshal([]byte(response.ToJsonString()), rspBody); err != nil {
		return nil, fmt.Errorf("chatByOnce-Unmarshal fail, rspBody: %s, err: %s", response.ToJsonString(), err)
	}

	return rspBody.Response, nil
}

// HyByStream 流式调用混元接口
func (p *HyCli) chatByStream(request *hunyuan.ChatCompletionsRequest,
	rspChan chan string, gapTime time.Duration, retryCount int) error {
	response, err := p.clientChat(request)
	if err != nil {
		close(rspChan)
		return err
	}

	closeCount := 0
	curRsp := ""
	ticker := time.NewTicker(gapTime)
	defer ticker.Stop()
	for {
		select {
		case event := <-response.Events:
			eventBody := &hyentity.HyChatRsp{}
			if err := json.Unmarshal(event.Data, eventBody); err != nil {
				continue
			}
			for _, choice := range eventBody.Choices {
				curRsp += choice.Delta.Content
			}

		case <-ticker.C:
			if closeCount == retryCount {
				close(rspChan)
				return nil
			}
			if len(curRsp) == 0 {
				closeCount++
			} else {
				closeCount = 0
			}

			rspChan <- curRsp
			curRsp = ""
		}
	}
}

func (p *HyCli) clientChat(request *hunyuan.ChatCompletionsRequest) (*hunyuan.ChatCompletionsResponse, error) {
	credential := common.NewCredential(p.secretID, p.secretKey)
	cpf := profile.NewClientProfile()
	request.Model = common.StringPtr(p.model)

	client, _ := hunyuan.NewClient(credential, regions.Guangzhou, cpf)
	response, err := client.ChatCompletions(request)
	util.WriteFile("req.txt", []string{request.ToJsonString()})
	util.WriteFile("rsp.txt", []string{response.ToJsonString()})
	if err != nil {
		return nil, fmt.Errorf("clientChat fail, err: %s", err)
	}
	return response, nil
}
