package services

import (
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func Get_process_id_from_process_name() int {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/fo", "csv", "/nh")
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, processus+winexe) {
				fields := strings.Split(line, ",")
				fields[1] = strings.Replace(fields[1], "\"", "", -1)
				id, err := strconv.Atoi(fields[1])
				if err != nil {
					log.Fatal(err)
				}
				return id
			}
		}
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("ps", "-ef")
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, processus+linexe) {
				fields := strings.Fields(line)
				id, err := strconv.Atoi(fields[1])
				if err != nil {
					log.Fatal(err)
				}
				return id
			}
		}
	}
	return -1
}
