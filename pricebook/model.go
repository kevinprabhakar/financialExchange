package pricebook

import (
	"gopkg.in/mgo.v2/bson"
	"financialExchange/api"
	"time"
)

//taken from GDAX API
//
//{
//"trade_id": 4729088,
//"price": "333.99",
//"size": "0.193",
//"bid": "333.98",
//"ask": "333.99",
//"volume": "5957.11914015",
//"time": "2015-11-14T20:46:03.511254Z"
//}

type PricePoint struct{
	Id				bson.ObjectId		`json:"id" bson:"_id"`
	TransactionId 	bson.ObjectId		`json:"transactionId" bson:"transactionId"`
	Price 			api.Money			`json:"price" bson:"price"`
	Size 			float64 			`json:"size" bson:"size"`
	Bid 			api.Money 			`json:"bid" bson:"bid"`
	Ask 			api.Money 			`json:"ask" bson:"ask"`
	Volume 			float64 			`json:"volume" bson:"volume"`
	time 			time.Time 			`json:"time" bson:"time"`
}
