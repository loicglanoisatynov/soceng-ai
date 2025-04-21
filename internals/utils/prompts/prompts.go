package prompts

import (
	env "soceng-ai/internals/server/env"
	colors "soceng-ai/internals/utils/colors"
	"time"
)

var (
	Error   = colors.Red + "[ERROR]  " + colors.Reset
	Log     = colors.Cyan + "[LOG]     " + colors.Reset
	Success = colors.Green + "[SUCCESS] " + colors.Reset
	Info    = colors.Yellow + "[INFO]   " + colors.Reset
	Debug   = colors.Purple + "[DEBUG]  " + colors.Reset
	Prompt  = colors.Cyan + "[ " + time.Now().Format("15:04:05.000000") + " server" + env.Get_dev_mode_as_string() + " ] " + colors.Reset
)
