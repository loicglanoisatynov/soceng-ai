package services

import (
	"os"
	"syscall"
)

const PID_FILE = "/tmp/socengai.pid"

const UID = 0
const GUID = 0
const default_host = "localhost"
const default_port = "80"

// Là où se trouve le serveur. Pointe vers le fichier exécutable du serveur (./server.exe)
var processus string = "socengai-server"
var winexe string = ".exe" // Extension pour les fichiers exécutables Windows
var linexe string = ""     // Extension pour les fichiers exécutables Linux
var binpath string = "./bin/"

/* Vérifie si le serveur est déjà en train de tourner. */
func check_if_running() bool {
	id := get_process_id_from_process_name()
	proc, err := os.FindProcess(id)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
