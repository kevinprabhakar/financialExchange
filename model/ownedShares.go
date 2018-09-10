package model

type OwnedShare struct{
	UserID 		int64 			`json:"userId"`
	Security 	int64 			`json:"security"`
	NumShares 	int 			`json:"numShares"`
}
