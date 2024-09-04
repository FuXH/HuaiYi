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
	role          *roleentity.Role
	promptTplFile string
	inputTpl      string
	llmConfig     *hyentity.HyChatConfig
}

func init() {
	translator := &Translator{
		promptTplFile: "./prompt/translator.tpl",
		llmConfig:     hyentity.NewChatConfig(),
	}
	translator.ParsePromptFile()

	role.RegisterRole(translator.role.Name, translator)
}

func (p *Translator) Do(input string) error {
	input = p.editInput(input)

	// 1、调用llm进行翻译
	chatTask := chat_task.Init(p.role, hunyuan.GetInstance(), p.llmConfig, nil)
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

func (p *Translator) Output(hyRsp *hyentity.HyChatRsp) error {
	//hyRsp.Display()
	content := hyRsp.GetContent(p.llmConfig.IsStream)
	fmt.Printf("翻译结果：%s\n", content)
	return nil
}

func (p *Translator) ParsePromptFile() {
	prompt, err := util.ReadFile(p.promptTplFile)
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

	p.role = &roleentity.Role{
		Name: transPpt.RoleName,
		Desc: transPpt.RoleDesc,
	}
	p.inputTpl = transPpt.InputTemplate
}

func (p *Translator) editInput(input string) string {
	input = strings.Replace(p.inputTpl, ReplaceQues, input, -1)
	return input
}
