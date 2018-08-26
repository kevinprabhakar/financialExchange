package entity

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entity struct{
	Id 			bson.ObjectId		`json:"id" bson:"_id"`
	Name		string 				`json:"name" bson:"name"`
	Email 		string 				`json:"email" bson:"email"`
	PassHash	string 				`json:"passHash" bson:"passHash"`
	Security 	bson.ObjectId 		`json:"security" bson:"security"`
	Created 	time.Time			`json:"created" bson:"created"`
	Deleted 	time.Time 			`json:"deleted" bson:"deleted"`
}

//From Client
type CreateEntityParams struct{
	Name 			string 				`json:"name"`
	Email 			string 				`json:"email"`
	Password 		string 				`json:"password"`
	PasswordVerify  string				`json:"passwordVerify"`
	Symbol 			string 				`json:"symbol"`
}