package cmds

import (
	"bigv.io/client/cmds/util"
	"bigv.io/client/util/log"
)

func (cmds *CommandSet) HelpForResize() util.ExitCode {
	log.Log("bigv resize")
	log.Log("")
	log.Log("usage: bigv resize disc <virtual machine> <disc> <size>")
	return util.E_USAGE_DISPLAYED
}

func (cmds *CommandSet) ResizeDisk(args []string) util.ExitCode {
	flags := util.MakeCommonFlagSet()
	flags.Parse(args)
	args = cmds.config.ImportFlags(flags)

	nameStr, ok := util.ShiftArgument(&args, "virtual machine")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	name := cmds.bigv.ParseGroupName(nameStr)

	discId, ok := util.ShiftArgument(&args, "disc id")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	sizeStr, ok := util.ShiftArgument(&args, "")
	if !ok {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	disc, err := strconv.parseInt(discId, 10, 32)
	if err != nil {
		cmds.HelpForList()
		return util.E_PEBKAC
	}

	size, err := util.ParseSize(sizeStr)
	if err != nil {
		return util.ProcessError(err)
	}

	err := cmds.EnsureAuth()
	if err != nil {
		return util.ProcessError(err)
	}

	cmds.bigv.ResizeDisc(name, disc, size)

	return util.E_SUCCESS
}
