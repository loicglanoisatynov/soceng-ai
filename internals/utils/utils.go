package utils

import "os"

func We_are_on_WSL() bool {
	return os.Getenv("WSL_DISTRO_NAME") != ""
}
