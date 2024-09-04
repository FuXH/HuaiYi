package tool_function

type Function interface {
	GetInfo() (string, string, string)
	Call(args string) string
}

var (
	FunctionList = make(map[string]Function)
)

func RegisterFunction(fcName string, function Function) {
	FunctionList[fcName] = function
}
