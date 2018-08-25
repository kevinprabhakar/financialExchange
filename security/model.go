package security

import (
	"gopkg.in/mgo.v2/bson"
	"financialExchange/pricebook"
)

type Security struct{
	Id 			bson.ObjectId 		`json:"id" bson:"_id"`
	Entity 		bson.ObjectId 		`json:"entity" bson:"entity"`
	Symbol 		string 				`json:"symbol" bson:"symbol"`
	PriceBook	bson.ObjectId 		`json:"priceBook" bson:"priceBook"`
}



