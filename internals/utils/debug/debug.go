package debug

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}