package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jvinaya/goapp/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {

	config, err := utils.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load Config : ", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("connot connect to database ", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
