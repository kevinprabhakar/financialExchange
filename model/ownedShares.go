package model

type OwnedShare struct{
	Id 			int64 			`json:"id"`
	UserID 		int64 			`json:"userId"`
	Security 	int64 			`json:"security"`
	//Implement BOOL isSecurity
	NumShares 	int 			`json:"numShares"`
}
