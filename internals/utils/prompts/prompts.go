package prompts

import (
	"fmt"
	env "soceng-ai/internals/server/env"
	colors "soceng-ai/internals/utils/colors"
	"time"
)

var (
	Error   = colors.Cyan + "| [" + colors.Red + " ERROR " + colors.Cyan + "] | " + colors.Reset
	Log     = colors.Cyan + "| [" + colors.Cyan + " LOG    " + colors.Cyan + "] | " + colors.Reset
	Success = colors.Cyan + "| [" + colors.Green + " SUCCESS " + colors.Cyan + "] | " + colors.Reset
	Info    = colors.Cyan + "| [" + colors.Yellow + " INFO   " + colors.Cyan + "] | " + colors.Reset
	Debug   = colors.Cyan + "| [" + colors.Purple + " DEBUG  " + colors.Cyan + "] | " + colors.Reset
	Warning = colors.Cyan + "| [" + colors.Purple + " WARNING " + colors.Cyan + "] | " + colors.Reset
	// Prompt_server = colors.Cyan + "[ " + time.Now().Format("15:04:05.000000") + " server" + env.Get_dev_mode_as_string() + " ] " + colors.Reset
	// Prompt_tests  = colors.Cyan + "[ " + time.Now().Format("15:04:05.000000") + " tests" + env.Get_dev_mode_as_string() + " ] " + colors.Reset
)

func Prompts_tests(time time.Time, payload string) {
	fmt.Println(colors.Cyan + "[ " + time.Format("15:04:05.000000") + " tests" + env.Get_dev_mode_as_string() + " ] " + colors.Reset + payload)
}

func Prompts_server(time time.Time, payload string) {
	fmt.Println(colors.Cyan + "[ " + time.Format("15:04:05.000000") + " server" + env.Get_dev_mode_as_string() + " ] " + colors.Reset + payload)
}
