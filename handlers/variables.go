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
	"financialExchange/transaction"
	"fmt"
)

type AccessToken struct{
	AccessToken 		string 		`json:"accessToken"`
}

var ServerLogger *util.Logger
var CustomerController *customer.CustomerController
var EntityController *entity.EntityController
var OrderController *order.OrderController
var TransactionController *transaction.TransactionController
var TransactionChannel chan model.OrderTransactionPackage
var ErrorChannel chan error

func init(){
	DBConn, err := sql.OpenDatabase("exchange")

	if err != nil{
		ServerLogger.Debug(err.Error())
		return
	}

	TransactionChannel = make(chan model.OrderTransactionPackage, 100)
	ErrorChannel = make(chan error, 100)
	ServerLogger = util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)


	CustomerController = customer.NewCustomerController(ServerLogger, DBConn)
	EntityController = entity.NewEntityController(DBConn, ServerLogger)
	OrderController = order.NewOrderController(DBConn, ServerLogger, TransactionChannel)
	TransactionController = transaction.NewTransactionController(DBConn, ServerLogger, TransactionChannel, ErrorChannel)


	//Database Table Creation
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
		return
	}

	err = DBConn.CreateTransactionTable()
	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	//Go Routines

	//Handling Matched Orders (Runs indefinitely)
	go TransactionController.HandleMatchedOrders()

	go func(errorChan chan error){
		for err := range errorChan{
			ServerLogger.ErrorMsg(fmt.Sprintf("Error while matching transactions: %s", err.Error()))
			return
		}
	}(ErrorChannel)
}


