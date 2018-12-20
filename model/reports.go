package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


type SecurityReport struct{
	Id 			bson.ObjectId 		`json:"id"`
	Entity 		 		`json:"entity"`
	Symbol 		string 				`json:"symbol"`
}

type MostTradedReport struct{
	Security 		Security 		`json:"security"`
	Frequency 		int64			`json:"frequency"`
}

type TransactionReport struct{
	Id			bson.ObjectId			`json:"id"`
	Buyer 		CustomerReport	`json:"buyer"`
	Seller 		CustomerReport `json:"seller"`
	Order 		OrderReport		`json:"order"`
	TotalCost	Money				`json:"totalCost"`
	SystemFee	Money				`json:"systemFee"`
	Created 	time.Time				`json:"created"`
	Executed 	time.Time				`json:"executed"`
	Status 		CompletionStatus	`json:"status"`
}

type PortfolioReport struct{
	Id 					bson.ObjectId						`json:"id"`
	User 				bson.ObjectId 			`json:"user"`
	Value 				Money 							`json:"value"`
	StockValue 			Money							`json:"stockValue"`
	CashValue			Money							`json:"cashValue"`
	WithdrawableFunds	Money 							`json:"withdrawableFunds"`
	OwnedShares			map[string]int			`json:"ownedShares"`
	Orders 				[]OrderReport						`json:"orders"`
	Transactions 		[]TransactionReport			`json:"transactions"`
}

type OrderReport struct{
	Id				bson.ObjectId				`json:"id"`
	Investor 		CustomerReport		`json:"investor"`
	Security 		SecurityReport		`json:"security"`
	Symbol 			string						`json:"symbol"`
	Action 			InvestorAction 			`json:"investorAction"`
	OrderType 		OrderType				`json:"orderType"`
	CostPerShare	Money					`json:"costPerShare"`
	CostOfShares	Money 					`json:"costOfShares"`
	SystemFee 		Money 					`json:"systemFee"`
	TotalCost 		Money					`json:"totalCost"`
	Created 		time.Time					`json:"created"`
	Updated 		time.Time 					`json:"updated"`
	Fulfilled		time.Time 					`json:"fulfilled"`
	Status 			CompletionStatus		`json:"orderStatus"`
	Transactions 	[]bson.ObjectId				`json:"transactions"`

	//Limit Orders Only
	AllowTakers		bool 						`json:"allowTakers"`
	LimitPerShare	Money 					`json:"limitPerShare"`
	TakerFee		Money 					`json:"takerFee"`

	//Stop Orders Only
	StopPrice 		Money 					`json:"stopPrice"`
	StopLimitPrice	Money 					`json:"stopLimitPrice"`
}

type EntityReport struct{
	Id 			bson.ObjectId					`json:"id"`
	Name		string 							`json:"name"`
	Email 		string 							`json:"email"`
	Security 	SecurityReport	 				`json:"security"`
}

type CustomerReport struct{
	Id 			bson.ObjectId		`json:"id"`
	FirstName	string 				`json:"firstName"`
	LastName 	string 				`json:"lastName"`
	Email 		string 				`json:"email"`
	Portfolio   PortfolioReport		`json:"portfolio"`
}

type PriceBookReport struct{

}