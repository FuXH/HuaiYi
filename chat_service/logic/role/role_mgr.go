package role

import (
	hyentity "chat_service/entity/hunyuan_msg"
	roleentity "chat_service/entity/role"
)

type Role interface {
	ParsePromptFile()
	Do(input string) error
	Input(input string) string
	Output(hyRsp *hyentity.HyChatRsp) error
}

type BaseRole struct {
	Role          *roleentity.Role
	LlmConfig     *hyentity.HyChatConfig
	FuncCallList  []*hyentity.HyTool
	PromptTplFile string
}

var (
	RoleMap = make(map[string]Role)
)

func RegisterRole(name string, role Role) {
	RoleMap[name] = role
}

func GetRole(name string) Role {
	return RoleMap[name]
}
