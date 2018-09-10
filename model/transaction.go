package model

import (

	"time"

)

type Transaction struct{
	OrderPlaced 		int64 				`json:"orderPlaced"`
	OrdersFulfilled 	[]int64 			`json:"ordersFulfilled"`
	TotalCost			Money				`json:"totalCost"`
	SystemFee			Money				`json:"systemFee"`
	Created 			time.Time			`json:"created"`
}

