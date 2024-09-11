package tool_function

import hyentity "chat_service/entity/hunyuan_msg"

type Function interface {
	GetInfo() (string, string, string)
	Call(args string) string
}

func ConvertHyTool(function Function) *hyentity.HyTool {
	fcName, fcDesc, argsDesc := function.GetInfo()
	return &hyentity.HyTool{
		Type: "function",
		Function: &hyentity.HyFunction{
			Name:        fcName,
			Parameters:  argsDesc,
			Description: fcDesc,
		},
	}
}

var (
	FunctionList = make(map[string]Function)
)

func RegisterFunction(fcName string, function Function) {
	FunctionList[fcName] = function
}

func GetFuncCallList(funcName ...string) []*hyentity.HyTool {
	res := make([]*hyentity.HyTool, 0)
	for _, fcName := range funcName {
		function, ok := FunctionList[fcName]
		if ok {
			res = append(res, ConvertHyTool(function))
		}
	}
	return res
}

func CallFunction(funcName string, args string) string {
	function, ok := FunctionList[funcName]
	if ok {
		return function.Call(args)
	}
	return ""
}
