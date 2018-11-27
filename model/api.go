package model

import "github.com/shopspring/decimal"

type Money struct{
	decimal.Decimal
}

type CompletionStatus int
const(
	Untouched 	CompletionStatus = 0
	InProgress	CompletionStatus = 1
	Finished  	CompletionStatus = 2
	Cancelled	CompletionStatus = 3
)

type InvestorAction int
const(
	Buy		InvestorAction = 0
	Sell 	InvestorAction = 1
)

type InvestorType int
const(
	Maker 	InvestorType = 0
	Taker 	InvestorType = 1

)

type OrderType int
const(
	MARKET		OrderType = 0
	LIMIT	 	OrderType = 1
	STOP		OrderType = 2
)

func NewMoneyObject(amount float64) (Money){
	quantity := decimal.NewFromFloat(amount)
	return Money{quantity}
}

func NewMoneyObjectFromDecimal(amount decimal.Decimal)(Money){
	return Money{amount}
}

type SearchParams struct{
	Prefix 		string 		`json:"prefix"`
}

type SearchResult struct{
	EntityName 	string 		`json:"entityName"`
	Symbol		string 		`json:"symbol"`
	SecurityID  int64 		`json:"securityID"`
}