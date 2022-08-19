package command

const (
	CMD_START_CLIENT Command = 1
	CMD_START_NN     Command = 2
	CMD_START_DN     Command = 3

	CMD_CLIENT_PUT    Command = 4
	CMD_CLIENT_GET    Command = 5
	CMD_CLIENT_DELETE Command = 6
	CMD_CLIENT_STAT   Command = 7
	CMD_CLIENT_RENAME Command = 8
	CMD_CLIENT_MKDIR  Command = 9
	CMD_CLIENT_LIST   Command = 10
)
