package colors

const (
	Red    = "\033[31m"
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"
	Yellow = "\033[33m"
)

func Yellow_ify(text string) string {
	return Yellow + text + Reset
}

func Cyan_ify(text string) string {
	return Cyan + text + Reset
}

func Red_ify(text string) string {
	return Red + text + Reset
}

func Green_ify(text string) string {
	return Green + text + Reset
}

func Purple_ify(text string) string {
	return Purple + text + Reset
}
