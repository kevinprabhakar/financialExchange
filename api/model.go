package api

import "github.com/shopspring/decimal"

type Money struct{
	decimal.Decimal
}

type OrderStatus int
const(
	Untouched 	OrderStatus = 0
	InProgress	OrderStatus = 1
	Finished  	OrderStatus = 2
	Cancelled	OrderStatus = 3
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

func NewMoneyObject(amount float64) (*Money){
	quantity := decimal.NewFromFloat(amount)
	return &Money{quantity}
}