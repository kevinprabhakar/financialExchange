package testing

import (
	"financialExchange/sql"
	"io/ioutil"
	"os"
	"financialExchange/util"
)

var TestingController *TestController

func init(){
	ServerLogger := util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	DBConn, err := sql.OpenDatabase("exchange")

	if err != nil{
		ServerLogger.Debug(err.Error())
		return
	}

	TestingController = NewTestController(DBConn, ServerLogger)

}
