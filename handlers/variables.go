package handlers

import (
	"io/ioutil"
	"os"
	"financialExchange/mongo"
	"financialExchange/util"
	"financialExchange/customer"
	"financialExchange/entity"
)

var MongoSession = mongo.GetMongoSession()
var ServerLogger = util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
var CustomerController = customer.NewCustomerController(MongoSession, ServerLogger)
var EntityController = entity.NewEntityController(MongoSession, ServerLogger)

type AccessToken struct{
	AccessToken 		string 		`json:"accessToken"`
}