package homeassistant

type DaliCommand interface {
	ReadReply(reply []byte) (err error)
	RequiresReply() bool
	DeviceType() byte
	Instruction() uint16
	SendTwice() bool
	Val() byte
}

type Command struct {
	cmdVal byte

	instruction uint16

	val        byte
	deviceType byte
	sendTwice  bool
	isReply    bool
}

func (cmd *Command) Instruction() uint16 {
	return cmd.instruction
}

func (cmd *Command) RequiresReply() bool {
	return cmd.isReply
}

func (cmd *Command) DeviceType() byte {
	return cmd.deviceType
}
func (cmd *Command) SendTwice() bool {
	return cmd.sendTwice
}
func (cmd *Command) Val() byte {
	return cmd.val
}
func (cmd *Command) ReadReply(reply []byte) (err error) {
	l := reply[3]
	_ = l
	cmd.val = reply[4]
	return nil
}

type DAPCCommand struct {
	Command
}

func NewDAPCCommand(addr Address, level int) *DAPCCommand {
	return &DAPCCommand{
		Command: Command{
			instruction: uint16(addr) | uint16(level),
		},
	}
}

type SpecialCommand struct {
	Command
}

func NewSpecialCommand(cmdVal byte, args ...byte) *SpecialCommand {
	specialCommand := &SpecialCommand{
		Command: Command{
			cmdVal:      cmdVal,
			instruction: uint16(cmdVal) << 8,
		},
	}
	if len(args) == 1 {
		specialCommand.instruction |= uint16(args[0])
	}
	return specialCommand
}

type StandardCommand struct {
	Command
}

func NewStandardCommand(addr Address, cmdVal byte, args ...byte) StandardCommand {
	standardCommand := StandardCommand{
		Command: Command{
			cmdVal:      cmdVal,
			instruction: uint16(addr) | 0x100 | uint16(cmdVal),
			isReply:     true,
		},
	}
	if len(args) == 1 {
		standardCommand.instruction |= uint16(args[0])
	}
	return standardCommand
}
