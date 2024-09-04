package alfred

import (
	roleentity "chat_service/entity/role"
	"chat_service/tool_function"
	"chat_service/tool_function/weather"
	"chat_service/util"
	"encoding/json"
	"fmt"

	hyentity "chat_service/entity/hunyuan_msg"
	"chat_service/logic/role"
)

type Alfred struct {
	//base_role.BaseRole

	role          *roleentity.Role
	promptTplFile string
	funcCallList  []*hyentity.HyTool
}

func init() {
	alfred := &Alfred{
		promptTplFile: "./prompt/alfred.tpl",
		funcCallList: []*hyentity.HyTool{
			hyentity.NewHyTool(tool_function.FunctionList[weather.FuncName]),
		},
	}
	alfred.ParsePromptFile()

	role.RegisterRole(alfred.role.Name, alfred)
}

func (p *Alfred) ParsePromptFile() {
	prompt, err := util.ReadFile(p.promptTplFile)
	if err != nil {
		panic(err)
	}
	type alfredPrompt struct {
		RoleName string `json:"role_name"`
		RoleDesc string `json:"role_desc"`
	}
	alfredPmt := &alfredPrompt{}
	if err = json.Unmarshal([]byte(prompt), alfredPmt); err != nil {
		panic(err)
	}

	p.role = &roleentity.Role{
		Name: alfredPmt.RoleName,
		Desc: alfredPmt.RoleDesc,
	}
}

func (p *Alfred) Do(input string) error {

}

func (p *Alfred) Chat(chatID string,
	msg *hyentity.HyMessage,
	chatCfg *hyentity.HyChatConfig) (*hyentity.HyChatRsp, error) {
	hyRsp, err := p.BaseRole.Chat(chatID, msg, chatCfg)
	if err != nil {
		return nil, err
	}
	_ = p.Output(chatID, hyRsp, chatCfg)
	return hyRsp, nil
}

func (p *Alfred) Output(chatID string, hyRsp *hyentity.HyChatRsp, chatCfg *hyentity.HyChatConfig) error {
	content := hyRsp.GetContent(chatCfg.IsStream)
	fmt.Println("Alfred output: ", content)
	return nil
}
