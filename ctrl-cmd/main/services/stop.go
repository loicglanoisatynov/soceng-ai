package services

import (
	"fmt"
	"os"
)

func Stop() {
	pid := Get_process_id_from_process_name()
	if pid == -1 {
		fmt.Println("No process to stop")
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Error while stopping process")
		return
	}
	err = process.Kill()
	if err != nil {
		fmt.Println("Error while stopping process")
		return
	}
}
