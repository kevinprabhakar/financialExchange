package order

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"eSportsExchange/api"
	"eSportsExchange/transaction"
)

type Order struct{
	Id				bson.ObjectId				`bson:"_id"`
	Investor 		bson.ObjectId				`bson:"investor"`
	Symbol 			bson.ObjectId				`bson:"symbol"`
	Action 			api.InvestorAction 			`bson:"investorAction"`
	OrderType 		api.OrderType				`bson:"orderType"`
	CostPerShare	api.Money					`bson:"costPerShare"`
	CostOfShares	api.Money 					`bson:"costOfShares"`
	SystemFee 		api.Money 					`bson:"systemFee"`
	Created 		time.Time					`bson:"created"`
	Updated 		time.Time 					`bson:"updated"`
	Fulfilled		time.Time 					`bson:"fulfilled"`
	Status 			api.OrderStatus				`bson:"orderStatus"`
	Transactions 	[]transaction.Transaction	`bson:"transactions"`

	//Limit Orders Only
	AllowTakers		bool 						`bson:"allowTakers"`
	LimitPerShare	api.Money 					`bson:"limitPerShare"`

	//todo: Look at Gdax Limit Orders to find out what TimeInForce does

	//Stop Orders Only
	StopPrice 		api.Money 					`bson:"stopPrice"`
	StopLimitPrice	api.Money 					`bson:"stopLimitPrice"`
}
