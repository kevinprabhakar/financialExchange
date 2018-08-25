package customer

import (
	"gopkg.in/mgo.v2/bson"
)

type Customer struct{
	Id 			bson.ObjectId		`json:"id" bson:"_id"`
	FirstName	string 				`json:"firstName" bson:"firstName"`
	LastName 	string 				`json:"lastName" bson:"lastName"`
	Email 		string 				`json:"email" bson:"email"`
	PassHash	string 				`json:"passHash" bson:"passHash"`
	Portfolio   bson.ObjectId		`json:"portfolio" bson:"portfolio"`
}

type CustomerSignUpParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
	PasswordVerify	string 			`json:"passwordVerify"`
	FirstName 		string 			`json:"firstName"`
	LastName 		string 			`json:"lastName"`
}

type CustomerSignInParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
}