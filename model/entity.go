package model

type Entity struct{
	Name		string 				`json:"name"`
	Email 		string 				`json:"email"`
	PassHash	string 				`json:"passHash"`
	Security 	int64					`json:"security"`
	Created 	int64			`json:"created"`
	Deleted 	int64 			`json:"deleted"`
}

//From Client
type CreateEntityParams struct{
	Name 			string 				`json:"name"`
	Email 			string 				`json:"email"`
	Password 		string 				`json:"password"`
	PasswordVerify  string				`json:"passwordVerify"`
	CreateSecurity 	bool 				`json:"createSecurity"`
	Symbol 			string 				`json:"symbol"`
}

type SignInEntityParams struct{
	Email 			string 				`json:"email"`
	Password 		string 				`json:"password"`
}