package main

import (
	"flag"
	"fmt"
	"os"

	"agit.com/smartdashboard-backend/config"
	"agit.com/smartdashboard-backend/helper"
	"agit.com/smartdashboard-backend/router"
)

// Init our logs and server
func Init() {
	helper.InitLogFile()
	helper.Log.Println("Smart Dashboard Init Start")
	config.PopulateDataInRev()
	config := config.GetConfig()
	r := router.Lists()
	r.Run(config.GetString("server.port"))
}

func main() {
	enviroment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*enviroment)
	config.LoadDB()
	config.MigrateDB()
	Init()
}
