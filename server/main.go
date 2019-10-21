package main

import "os"

func main() {
	listenAddr := ":5000"
	logPath := "./output.log"

	if len(os.Args) == 2 {
		listenAddr = os.Args[1]
	} else if len(os.Args) == 3 {
		listenAddr = os.Args[1]
		logPath = os.Args[2]
	}

	instantiateTemplates()
	serve(listenAddr, logPath)
}
