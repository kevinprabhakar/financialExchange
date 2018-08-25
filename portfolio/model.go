package portfolio

import (
	"gopkg.in/mgo.v2/bson"
	"financialExchange/api"
)

type Portfolio struct{
	Id 					bson.ObjectId				`json:"id" bson:"_id"`
	User 				bson.ObjectId 				`json:"user" bson:"user"`
	Value 				api.Money 					`json:"value" bson:"value"`
	StockValue 			api.Money					`json:"stockValue" bson:"stockValue"`
	CashValue			api.Money					`json:"cashValue" bson:"cashValue"`
	WithdrawableFunds	api.Money 					`json:"withdrawableFunds" bson:"withdrawableFunds"`
	OwnedShares			map[bson.ObjectId]int		`json:"ownedShares" bson:"ownedShares"`
	Orders 				[]bson.ObjectId				`json:"orders" bson:"orders"`
	Transactions 		[]bson.ObjectId 			`json:"transactions" bson:"transactions"`
}


