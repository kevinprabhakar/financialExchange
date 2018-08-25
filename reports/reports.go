package reports

import (
	"gopkg.in/mgo.v2/bson"
	"financialExchange/pricebook"
	"time"
	"financialExchange/api"
)


type SecurityReport struct{
	Id 			bson.ObjectId 		`json:"id" bson:"_id"`
	Entity 		bson.ObjectId 		`json:"entity" bson:"entity"`
	Symbol 		string 				`json:"symbol" bson:"symbol"`
	PriceBook	pricebook.PriceBook 		`json:"priceBook" bson:"priceBook"`
}

type TransactionReport struct{
	Id			bson.ObjectId			`json:"id" bson:"_id"`
	Buyer 		CustomerReport	`json:"buyer" bson:"buyer"`
	Seller 		CustomerReport `json:"seller" bson:"seller"`
	Order 		OrderReport		`json:"order" bson:"order"`
	TotalCost	api.Money				`json:"totalCost" bson:"totalCost"`
	SystemFee	api.Money				`json:"systemFee" bson:"systemFee"`
	Created 	time.Time				`json:"created" bson:"created"`
	Executed 	time.Time				`json:"executed" bson:"executed"`
	Status 		api.CompletionStatus	`json:"status" bson:"status"`
}

type PortfolioReport struct{
	Id 					bson.ObjectId						`json:"id" bson:"_id"`
	User 				bson.ObjectId 			`json:"user" bson:"user"`
	Value 				api.Money 							`json:"value" bson:"value"`
	StockValue 			api.Money							`json:"stockValue" bson:"stockValue"`
	CashValue			api.Money							`json:"cashValue" bson:"cashValue"`
	WithdrawableFunds	api.Money 							`json:"withdrawableFunds" bson:"withdrawableFunds"`
	OwnedShares			map[string]int			`json:"ownedShares" bson:"ownedShares"`
	Orders 				[]OrderReport						`json:"orders" bson:"orders"`
	Transactions 		[]TransactionReport			`json:"transactions" bson:"transactions"`
}

type OrderReport struct{
	Id				bson.ObjectId				`json:"id" bson:"_id"`
	Investor 		CustomerReport		`json:"investor" bson:"investor"`
	Security 		SecurityReport		`json:"security" bson:"security"`
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

type EntityReport struct{
	Id 			bson.ObjectId					`json:"id" bson:"_id"`
	Name		string 							`json:"name" bson:"name"`
	Email 		string 							`json:"email" bson:"email"`
	Security 	SecurityReport	 				`json:"security" bson:"security"`
}

type CustomerReport struct{
	Id 			bson.ObjectId		`json:"id" bson:"_id"`
	FirstName	string 				`json:"firstName" bson:"firstName"`
	LastName 	string 				`json:"lastName" bson:"lastName"`
	Email 		string 				`json:"email" bson:"email"`
	Portfolio   PortfolioReport		`json:"portfolio" bson:"portfolio"`
}

type PriceBookReport struct{

}