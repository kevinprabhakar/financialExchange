package model

type OwnedShare struct{
	Id 			int64 			`json:"id"`
	UserID 		int64 			`json:"userId"`
	Security 	int64 			`json:"security"`
	//Implement BOOL isSecurity
	NumShares 	int 			`json:"numShares"`
}

type OwnedShareReport struct{
	UserID 		int64 			`json:"userId"`
	Security 	int64 			`json:"security"`
	Symbol 		string 			`json:"symbol"`
	CurrPrice 	float64			`json:"currPrice"`
	EntityName	string			`json:"entityName"`
	//Implement BOOL isSecurity
	NumShares 	int 			`json:"numShares"`
}