package main

import (
	"github.com/cucyber/cyberrange/services/webserver/db"
	"github.com/cucyber/cyberrange/services/webserver/server"
	"os"
)

func main() {
	listenAddr := ":5000"
	logPath := "./output.log"

	if len(os.Args) == 2 {
		listenAddr = os.Args[1]
	} else if len(os.Args) == 3 {
		listenAddr = os.Args[1]
		logPath = os.Args[2]
	}

	db.InitializeDatabase()
	defer db.CloseDatabase()

	server.Serve(listenAddr, logPath)
}
