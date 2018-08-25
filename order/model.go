package order

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"financialExchange/api"

)

type Order struct{
	Id				bson.ObjectId				`json:"id" bson:"_id"`
	Investor 		bson.ObjectId				`json:"investor" bson:"investor"`
	Security 		bson.ObjectId				`json:"security" bson:"security"`
	Symbol 			string						`json:"symbol" bson:"symbol"`
	Action 			api.InvestorAction 			`json:"investorAction" bson:"investorAction"`
	OrderType 		api.OrderType				`json:"orderType" bson:"orderType"`
	CostPerShare	api.Money					`json:"costPerShare" bson:"costPerShare"`
	CostOfShares	api.Money 					`json:"costOfShares" bson:"costOfShares"`
	SystemFee 		api.Money 					`json:"systemFee" bson:"systemFee"`
	TotalCost 		api.Money					`json:"totalCost" bson:"totalCost"`
	Created 		time.Time					`json:"created" bson:"created"`
	Updated 		time.Time 					`json:"updated" bson:"updated"`
	Fulfilled		time.Time 					`json:"fulfilled" bson:"fulfilled"`
	Status 			api.CompletionStatus		`json:"orderStatus" bson:"orderStatus"`
	Transactions 	[]bson.ObjectId				`json:"transactions" bson:"transactions"`

	//Limit Orders Only
	AllowTakers		bool 						`json:"allowTakers" bson:"allowTakers"`
	LimitPerShare	api.Money 					`json:"limitPerShare" bson:"limitPerShare"`
	TakerFee		api.Money 					`json:"takerFee" bson:"takerFee"`

	//Stop Orders Only
	StopPrice 		api.Money 					`json:"stopPrice" bson:"stopPrice"`
	StopLimitPrice	api.Money 					`json:"stopLimitPrice" bson:"stopLimitPrice"`
}


