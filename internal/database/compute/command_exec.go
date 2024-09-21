package compute

type CommandType int

var (
	CommandTypeGet CommandType = 1
	CommandTypeSet CommandType = 2
	CommandTypeDel CommandType = 3
)

type CommandExecResult struct {
	Result string
}

type CommandExec struct {
	Command CommandType
	Key     string
	Val     string
}
