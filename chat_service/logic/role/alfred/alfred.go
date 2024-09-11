package alfred

import (
	roleentity "chat_service/entity/role"
	"chat_service/logic/task/chat_task"
	"chat_service/repository/remote/hunyuan"
	"chat_service/repository/storage/tcvectordb"
	"chat_service/tool_function"
	"chat_service/tool_function/weather"
	"chat_service/util"
	"encoding/json"
	"fmt"

	hyentity "chat_service/entity/hunyuan_msg"
	"chat_service/logic/role"
)

type Alfred struct {
	*role.BaseRole

	db *tcvectordb.TCVectorDB // 记忆
}

func init() {
	alfred := &Alfred{
		BaseRole: &role.BaseRole{
			Role: &roleentity.Role{
				Name: "Alfred",
			},
			PromptTplFile: "./prompt/alfred.tpl",
			LlmConfig:     hyentity.NewChatConfig(),
			FuncCallList:  tool_function.GetFuncCallList(weather.FuncName),
		},
	}
	alfred.ParsePromptFile()

	role.RegisterRole(alfred.Role.Name, alfred)
}

func (p *Alfred) ParsePromptFile() {
	prompt, err := util.ReadFile(p.PromptTplFile)
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

	p.Role.Desc = alfredPmt.RoleDesc
}

func (p *Alfred) Do(input string) error {
	chatTask := chat_task.Init(p.Role, hunyuan.GetInstance(), p.LlmConfig, p.FuncCallList)
	chatTaskRsp, err := chatTask.Exec(input)
	if err != nil {
		return err
	}

	if err = p.Output(chatTaskRsp); err != nil {
		return err
	}
	return nil
}

func (p *Alfred) Input(input string) string {
	return input
}

func (p *Alfred) Output(hyRsp *hyentity.HyChatRsp) error {
	content := hyRsp.GetContent(p.LlmConfig.IsStream)
	fmt.Println("Alfred output: ", content)
	return nil
}

func (p *Alfred) Memory(input string) error {
	// 1、通过metadata过滤文件夹

	// 2、通过文本匹配过滤匹配度高的内容

	// 3、二者去重

	return nil
}
