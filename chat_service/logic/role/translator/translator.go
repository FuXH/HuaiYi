package translator

import (
	roleentity "chat_service/entity/role"
	"chat_service/logic/task"
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
}

func init() {
	translator := &Translator{
		promptTplFile: "./prompt/translator.tpl",
	}
	translator.ParsePromptFile()

	role.RegisterRole(translator.role.Name, translator)
}

func (p *Translator) Chat(chatID string,
	msg *hyentity.HyMessage, chatCfg *hyentity.HyChatConfig) (*hyentity.HyChatRsp, error) {
	msg.Content = p.EditMsg(msg)
	hyRsp, err := p.BaseRole.Chat(chatID, msg, chatCfg)
	if err != nil {
		return nil, err
	}

	_ = p.Output(chatID, hyRsp, chatCfg)
	return hyRsp, nil
}

func (p *Translator) Do(input string) error {
	chatTask := chat_task.Init(p.role, hunyuan.GetInstance(), hyentity.NewChatConfig(), nil)
	taskList := []task.Task{chatTask}

	input = p.editInput(input)
	for _, task := range taskList {
		taskRsp, err := task.Exec(input)
		
	}
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

func (p *Translator) Output(chatID string, hyRsp *hyentity.HyChatRsp, chatCfg *hyentity.HyChatConfig) error {
	//hyRsp.Display()
	content := hyRsp.GetContent(chatCfg.IsStream)
	fmt.Println(content)
	return nil
}
