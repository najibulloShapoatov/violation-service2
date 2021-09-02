package main

import (
	"fmt"
	"log"
	"os"
	"service/worker/app"
)

var (
	logFileName = "./log/worker_%s.log"
	confFile    = "./conf/default.conf"
)

func main() {

	app := app.New(logFileName, confFile, isRunning())

	app.Run()
}

/*
--------
--------
--------
--------
--------
--------
--------
--------
--------
*/

func isRunning() *os.File {

	fileStruct, err := os.Stat("processid")
	if err == nil {
		if fileStruct.Size() > 0 {
			log.Println("Application already running")
			os.Exit(0)
		}
	}
	return createPID()
}

func createPID() *os.File {
	pid := os.Getpid()
	processidFile, err := os.OpenFile("processid", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	processidFile.Write([]byte(fmt.Sprint(pid)))
	return processidFile
}
