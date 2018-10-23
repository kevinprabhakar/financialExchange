package model


type Order struct{
	ID 				int64 							`json:"id"`
	Investor 		int64							`json:"investor"`
	Security 		int64							`json:"security"`
	Symbol 			string						`json:"symbol"`
	InvestorAction 	InvestorAction 			`json:"investorAction"`
	InvestorType 	InvestorType 			`json:"investorType"`
	OrderType 		OrderType				`json:"orderType"`
	NumShares 		int							`json:"numShares"`
	NumSharesRemaining	int 				`json:"numSharesRemaining"`
	//Implement numSharesRemaining
	CostPerShare	Money					`json:"costPerShare"`
	CostOfShares	Money 					`json:"costOfShares"`
	SystemFee 		Money 					`json:"systemFee"`
	TotalCost 		Money					`json:"totalCost"`

	//UTC Time in Seconds
	Created 		int64						`json:"created"`
	Updated 		int64 					`json:"updated"`
	Fulfilled		int64 					`json:"fulfilled"`
	Status 			CompletionStatus		`json:"orderStatus"`

	//Removed transaction list from Orders

	//Limit Orders Only
	AllowTakers		bool 						`json:"allowTakers"`
	LimitPerShare	Money 					`json:"limitPerShare"`
	TakerFee		Money 					`json:"takerFee"`

	//Stop Orders Only
	StopPrice 		Money 					`json:"stopPrice"`
}


type OrderCreateParams struct{
	UserID 				int64 					`json:"userId"`
	InvestorAction 		int 					`json:"investorAction"`
	InvestorType 		int 					`json:"investorType"`
	OrderType 			int 					`json:"orderType"`
	Symbol 				string 					`json:"symbol"`
	NumShares		 	int						`json:"numShares"`
	CostPerShare 		float64 				`json:"costPerShare"`
	TimeCreated			int64					`json:"timeCreated"`

	//Advanced features
	AllowTakers 		bool 					`json:"allowTakers"`
	LimitPerShare		float64 				`json:"limitPerShare"`
	StopPrice 			float64 				`json:"stopPrice"`
}

type OrderId struct{
	Id 			int64 							`json:"orderID"`
}

type OrderTransactionPackage struct{
	MainOrder 			int64 					`json:"mainOrder"`
	MatchingOrders 		map[int64]int 			`json:"matchingOrders"`
}

type IPOParams struct{
	SharePrice 			float64 			`json:"sharePrice"`
	NumShares 			int 			`json:"numShares"`
}