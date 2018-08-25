package transaction

import (
	"financialExchange/api"
	"gopkg.in/mgo.v2/bson"
	"time"
	"financialExchange/customer"
	"financialExchange/order"
)

type Transaction struct{
	Id			bson.ObjectId			`json:"id" bson:"_id"`
	Buyer 		bson.ObjectId			`json:"buyer" bson:"buyer"`
	Seller 		bson.ObjectId 			`json:"seller" bson:"seller"`
	Order 		bson.ObjectId			`json:"order" bson:"order"`
	TotalCost	api.Money				`json:"totalCost" bson:"totalCost"`
	SystemFee	api.Money				`json:"systemFee" bson:"systemFee"`
	Created 	time.Time				`json:"created" bson:"created"`
	Executed 	time.Time				`json:"executed" bson:"executed"`
	Status 		api.CompletionStatus	`json:"status" bson:"status"`
}

