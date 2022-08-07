package main

import (
	"os"
	"runtime"

	"github.com/TremblingV5/CrazyDFS/command"
	"github.com/TremblingV5/CrazyDFS/utils"
)

func Init() {
	utils.InitConfig()
	utils.InitLogger()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	Init()

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			args = []string{"--version"}
		}
	}

	newArgs := append([]string{os.Args[0]}, args...)
	command.Parse(newArgs)
}
