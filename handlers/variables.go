package handlers

import (
	"io/ioutil"
	"os"
	"financialExchange/sql"
	"financialExchange/util"
	"financialExchange/customer"
	"financialExchange/entity"
	"financialExchange/order"
	"financialExchange/model"
)

type AccessToken struct{
	AccessToken 		string 		`json:"accessToken"`
}

var ServerLogger *util.Logger
var CustomerController *customer.CustomerController
var EntityController *entity.EntityController
var OrderController *order.OrderController
var TransactionChannel chan model.OrderTransactionPackage

func init(){
	DBConn, err := sql.OpenDatabase("exchange")

	if err != nil{
		ServerLogger.Debug(err.Error())
		return
	}

	TransactionChannel = make(chan model.OrderTransactionPackage, 100)

	CustomerController = customer.NewCustomerController(ServerLogger, DBConn)
	EntityController = entity.NewEntityController(DBConn, ServerLogger)
	OrderController = order.NewOrderController(DBConn, ServerLogger, TransactionChannel)

	ServerLogger = util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//defer DBConn.DB.Close()

	err = DBConn.CreateCustomerTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	err = DBConn.CreatePortfolioTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	err = DBConn.CreateEntityTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	err = DBConn.CreateSecurityTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	err = DBConn.CreateOrderTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	err = DBConn.CreateOwnedSharesTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
	}
}


