package env

var dev_mode bool

func SetDevMode(mode bool) {
	dev_mode = mode
}

func Get_dev_mode() bool {
	return dev_mode
}

func Get_dev_mode_as_string() string {
	if dev_mode {
		return " dev"
	} else {
		return ""
	}
}
