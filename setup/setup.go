package setup

import (
	"flag"
	"log"
	"os"

	"github.com/Borwe/go-lsp/logs"
)

func getLogFile() (bool, *string) {
	if len(os.Args) == 1 {
		return false, nil
	}
	//test mode for testing ./tests files
	if flag.Lookup("test.v") != nil {
		return false, nil
	}
	if os.Args[1] == "--stdio" {
		if len(os.Args) == 2 {
			return false, nil
		}
		return true, &os.Args[2]
	} else {
		return true, &os.Args[1]
	}
}

func SetupLog() {
	foundFile, file := getLogFile()
	if foundFile {
		f, err := os.Open(*file)
		if err != nil {
			f, err = os.Create(*file)
			if err != nil {
				foundFile = false
			}
		}
		logs.Init(log.New(f, "", 0))
		return
	}
	logs.Init(log.New(os.Stderr, "", 0))
}
