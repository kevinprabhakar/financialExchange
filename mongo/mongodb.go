package mongo

import (
	"gopkg.in/mgo.v2"
	"os"
	"fmt"
)

const DBname = "financialExchange"

var debug = true

func GetMongoSession() (*mgo.Session){
	if (!debug){
		uri := os.Getenv("MONGODB_URI")
		session, err := mgo.Dial(uri)

		if (err != nil){
			panic(err)
		}
		return session
	}else{
		uri := "mongodb://localhost:27017/financialExchange"
		session, err := mgo.Dial(uri)

		if (err != nil){
			panic(err)
		}
		return session
	}

}

func GetDataBase(session *mgo.Session) *mgo.Database{
	return session.DB(DBname)
}

func GetCustomerCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Customer")
}

func GetPortfolioCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Portfolio")
}

func GetEntityCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Entity")
}

func GetSecurityCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Security")
}

func GetOrderBookForSecurity(db *mgo.Database, symbol string) *mgo.Collection{
	return db.C(fmt.Sprintf("%s-orders",symbol))
}

func GetTransactionsForSecurity(db *mgo.Database, symbol string) *mgo.Collection{
	return db.C(fmt.Sprintf("%s-transactions",symbol))
}

func GetPriceBookForSecurity(db *mgo.Database, symbol string) *mgo.Collection{
	return db.C(fmt.Sprintf("%s-prices"))
}