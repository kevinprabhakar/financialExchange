package entity

import (
	"gopkg.in/mgo.v2/bson"
	"financialExchange/security"
)

type Entity struct{
	Id 			bson.ObjectId		`json:"id" bson:"_id"`
	Name		string 				`json:"name" bson:"name"`
	Email 		string 				`json:"email" bson:"email"`
	PassHash	string 				`json:"passHash" bson:"passHash"`
	Security 	bson.ObjectId 		`json:"security" bson:"security"`
}

