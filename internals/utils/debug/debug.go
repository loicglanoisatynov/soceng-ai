package debug

import (
	"log"
)

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}

func LogError(msg string, err error) {
	log.Printf("[ERROR] %s: %v", msg, err)
}
