package translator

import (
	roleentity "chat_service/entity/role"
	"chat_service/logic/task/chat_task"
	"chat_service/repository/remote/hunyuan"
	"encoding/json"
	"fmt"
	"strings"

	hyentity "chat_service/entity/hunyuan_msg"
	"chat_service/logic/role"
	"chat_service/util"
)

const (
	ReplaceQues = "<<<input>>>"
)

type Translator struct {
	*role.BaseRole
	inputTpl string // 输入模版
}

func init() {
	translator := &Translator{
		BaseRole: &role.BaseRole{
			PromptTplFile: "./prompt/translator.tpl",
			Role: &roleentity.Role{
				Name: "Translator",
			},
			LlmConfig:    hyentity.NewChatConfig(),
			FuncCallList: nil,
		},
	}
	translator.ParsePromptFile()

	role.RegisterRole(translator.Role.Name, translator)
}

func (p *Translator) Do(input string) error {
	input = p.Input(input)

	// 1、调用llm进行翻译
	chatTask := chat_task.Init(p.Role, hunyuan.GetInstance(), p.LlmConfig, nil)
	chatTaskRsp, err := chatTask.Exec(input)
	if err != nil {
		return err
	}

	// 2、输出结果
	if err = p.Output(chatTaskRsp); err != nil {
		return err
	}
	return nil
}

func (p *Translator) Input(input string) string {
	input = strings.Replace(p.inputTpl, ReplaceQues, input, -1)
	return input
}

func (p *Translator) Output(hyRsp *hyentity.HyChatRsp) error {
	//hyRsp.Display()
	content := hyRsp.GetContent(p.LlmConfig.IsStream)
	fmt.Printf("翻译结果：%s\n", content)
	return nil
}

func (p *Translator) ParsePromptFile() {
	prompt, err := util.ReadFile(p.PromptTplFile)
	if err != nil {
		panic(err)
	}
	type transPrompt struct {
		RoleName      string `json:"role_name"`
		RoleDesc      string `json:"role_desc"`
		InputTemplate string `json:"input_template"`
	}
	transPpt := &transPrompt{}
	if err := json.Unmarshal([]byte(prompt), transPpt); err != nil {
		panic(err)
	}

	p.Role.Desc = transPpt.RoleDesc
	p.inputTpl = transPpt.InputTemplate
}
