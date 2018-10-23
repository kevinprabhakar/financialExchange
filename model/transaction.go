package model

type Transaction struct{
	Id 					int64 				`json:"id"`
	OrderPlaced 		int64 				`json:"orderPlaced"`
	OrdersFulfilling 	[]int64 			`json:"ordersFulfilling"`
	NumShares			int64				`json:"numShares"`
	CostPerShare 		Money 				`json:"costPerShare"`
	TotalCost			Money				`json:"totalCost"`
	SystemFee			Money				`json:"systemFee"`
	Created 			int64				`json:"created"`
	Security 			int64 				`json:"security"`
}

