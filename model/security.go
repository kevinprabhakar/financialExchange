package model

type Security struct{
	Id 			int64 		`json:"id"`
	Entity 		int64 		`json:"entity"`
	Symbol 		string 				`json:"symbol"`
	Created 	int64 		`json:"created"`
}
