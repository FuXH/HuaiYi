package role

import hyentity "chat_service/entity/hunyuan_msg"

type Role interface {
	ParsePromptFile()
	Do(input string) error
	Output(hyRsp *hyentity.HyChatRsp) error
	//EditMsg(msg *hyentity.HyMessage) string
	//Chat(chatID string, msg *hyentity.HyMessage, chatCfg *hyentity.HyChatConfig) (*hyentity.HyChatRsp, error)
	//CallTool(chatID string, hyRsp *hyentity.HyChatRsp, chatCfg *hyentity.HyChatConfig) (*hyentity.HyChatRsp, error)
	//GenerateID() string
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
