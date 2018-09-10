package model

type Portfolio struct{
	Customer 			int64 						`json:"customer"`
	Value 				Money 					`json:"value"`
	StockValue 			Money					`json:"stockValue"`
	CashValue			Money					`json:"cashValue"`
	WithdrawableFunds	Money 					`json:"withdrawableFunds"`
}


