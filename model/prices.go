package model

type PriceRequest struct{
	TimePeriod 		int 		`json:"timePeriod"`
	Security 		int64 		`json:"security"`
}

type PricePoint struct{
	TimeStamp		int64 		`json:"timeStamp"`
	PricePoint 		float64 	`json:"pricePoint"`
}