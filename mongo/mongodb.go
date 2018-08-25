package mongo

import (
	"gopkg.in/mgo.v2"
	"os"
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
	return db.C("portfolio")
}

func GetEntityCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Entity")
}

func GetSecurityCollection(db *mgo.Database) *mgo.Collection {
	return db.C("Security")
}

func GetOrderBook(db *mgo.Database) *mgo.Collection{
	return db.C("OrderBook")
}

func GetTransactions(db *mgo.Database) *mgo.Collection{
	return db.C("Transaction")
}

func GetPriceBookForSecurity(db *mgo.Database, symbol string) *mgo.Collection{
	return db.C(symbol)
}