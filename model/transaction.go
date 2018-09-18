package model

type Transaction struct{
	OrderPlaced 		int64 				`json:"orderPlaced"`
	OrdersFulfilling 	[]int64 			`json:"ordersFulfilling"`
	TotalCost			Money				`json:"totalCost"`
	SystemFee			Money				`json:"systemFee"`
	Created 			int64				`json:"created"`
}

