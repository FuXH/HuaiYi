package alfred

import (
	roleentity "chat_service/entity/role"
	"chat_service/logic/task/chat_task"
	"chat_service/repository/remote/hunyuan"
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
	llmConfig     *hyentity.HyChatConfig
	funcCallList  []*hyentity.HyTool
}

func init() {
	alfred := &Alfred{
		promptTplFile: "./prompt/alfred.tpl",
		llmConfig:     hyentity.NewChatConfig(),
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
	chatTask := chat_task.Init(p.role, hunyuan.GetInstance(), p.llmConfig, p.funcCallList)
	chatTaskRsp, err := chatTask.Exec(input)
	if err != nil {
		return err
	}

	if err = p.Output(chatTaskRsp); err != nil {
		return err
	}
	return nil
}

func (p *Alfred) Output(hyRsp *hyentity.HyChatRsp) error {
	content := hyRsp.GetContent(p.llmConfig.IsStream)
	fmt.Println("Alfred output: ", content)
	return nil
}
