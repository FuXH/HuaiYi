package testcase_classify

import (
	hyentity "chat_service/entity/hunyuan_msg"
	roleentity "chat_service/entity/role"
	"chat_service/logic/role"
	"chat_service/logic/task/chat_task"
	"chat_service/repository/remote/hunyuan"
	"chat_service/util"
	"encoding/json"
)

type TestcaseType struct {
	Level1 string // 一级分类
	Level2 string // 二级分类
}

type TestcaseClassify struct {
	*role.BaseRole
}

func init() {
	classify := &TestcaseClassify{
		BaseRole: &role.BaseRole{
			PromptTplFile: "./prompt/testcase_classify.tpl",
			Role: &roleentity.Role{
				Name: "TestcaseClassify",
			},
			LlmConfig:    hyentity.NewChatConfig(),
			FuncCallList: nil,
		},
	}
	classify.ParsePromptFile()

	//role.RegisterRole("TestcaseClassify", classify)
}

func (p *TestcaseClassify) ParsePromptFile() {
	prompt, err := util.ReadFile(p.PromptTplFile)
	if err != nil {
		panic(err)
	}
	p.Role.Desc = prompt
}

func (p *TestcaseClassify) Do(input string) (*TestcaseType, error) {
	input = p.Input(input)

	chatTask := chat_task.Init(p.Role, hunyuan.GetInstance(), p.LlmConfig, nil)
	chatTaskRsp, err := chatTask.Exec(input)
	if err != nil {
		return nil, err
	}

	testcaseType, err := p.Output(chatTaskRsp)
	if err != nil {
		return nil, err
	}
	return testcaseType, nil
}

func (p *TestcaseClassify) Input(input string) string {
	return input
}

func (p *TestcaseClassify) Output(hyRsp *hyentity.HyChatRsp) (*TestcaseType, error) {
	content := hyRsp.GetContent(p.LlmConfig.IsStream)
	testcaseType := &TestcaseType{}
	if err := json.Unmarshal([]byte(content), testcaseType); err != nil {
		return nil, err
	}
	return testcaseType, nil
}
