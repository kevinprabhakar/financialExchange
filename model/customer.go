package model

type Customer struct{
	Id 			int64 				`json:"id"`
	FirstName	string 				`json:"firstName"`
	LastName 	string 				`json:"lastName"`
	Email 		string 				`json:"email"`
	PassHash	string 				`json:"passHash"`
	Portfolio   int64				`json:"portfolio"`

}

type CustomerSignUpParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
	PasswordVerify	string 			`json:"passwordVerify"`
	FirstName 		string 			`json:"firstName"`
	LastName 		string 			`json:"lastName"`
}

type CustomerSignInParams struct {
	Email 			string 			`json:"email"`
	Password 		string 			`json:"password"`
}

type CustomerGiveMoneyParams struct {
	MoneyIncrease 	float64 		`json:"moneyIncrease"`
}