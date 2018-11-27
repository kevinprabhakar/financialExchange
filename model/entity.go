package model

type Entity struct{
	Id 			int64 				`json:"id"`
	Name		string 				`json:"name"`
	Email 		string 				`json:"email"`
	PassHash	string 				`json:"passHash"`
	Security 	int64					`json:"security"`
	Created 	int64				`json:"created"`
	Deleted 	int64 				`json:"deleted"`

	//Indicates whether or not entity has IPO'd
	//0 = no, 1+ = Transaction ID of IPO
	IPO 		int64				`json:"ipo"`
	AssocUser	int64 			`json:"assocUser"`
}

//From Client
type CreateEntityParams struct{
	Name 			string 				`json:"name"`
	Email 			string 				`json:"email"`
	Password 		string 				`json:"password"`
	PasswordVerify  string				`json:"passwordVerify"`
	Symbol 			string 				`json:"symbol"`
	CreateSecurity 	bool 				`json:"createSecurity"`
}

type SignInEntityParams struct{
	Email 			string 				`json:"email"`
	Password 		string 				`json:"password"`
}