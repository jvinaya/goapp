package main

import (
	"log"

	apihandler "github.com/jvinaya/goapp/apiHandler"
	"github.com/jvinaya/goapp/db"
	"github.com/jvinaya/goapp/postgres"
	"github.com/jvinaya/goapp/utils"
	_ "github.com/lib/pq"
)

func main() {

	config, err := utils.LoadConfig("./app.env")
	if err != nil {
		log.Fatal("Cannot load Config : ", err)
	}
	// Init database
	postgres.InitDatabase()
	store := db.NewStore(conn)
	server, err := apihandler.NewServer(config, store)
	if err != nil {
		log.Fatal("Connot Create  Server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Connot Connect to Server:", err)
	}

}
