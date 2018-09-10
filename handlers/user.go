package handlers

import (
	"net/http"
	"encoding/json"
	"financialExchange/util"
	"fmt"
	"financialExchange/model"
)

func SignUpCustomer(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.CustomerSignUpParams
	err := decoder.Decode(&params)

	if err != nil{
		http.Error(w, err.Error(), 400)
		return
	}
	ServerLogger.Debug(fmt.Sprintf("Signing Up %s", params.Email))
	err = CustomerController.SignUp(params)

	if err != nil{
		http.Error(w, err.Error(), 500)
		ServerLogger.ErrorMsg(fmt.Sprintf("Error while signing up %s: %s", params.Email, err.Error()))
		return
	}

	fmt.Fprint(w, util.GetNoDataSuccessResponse())
}

func SignInCustomer(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.CustomerSignInParams
	err := decoder.Decode(&params)

	if err != nil{
		http.Error(w, err.Error(), 400)
		return
	}

	ServerLogger.Debug(fmt.Sprintf("Signing In %s", params.Email))

	accessToken, err := CustomerController.SignIn(params)

	if err != nil{
		http.Error(w, err.Error(), 401)
		ServerLogger.ErrorMsg(fmt.Sprintf("Error while signing in %s: %s", params.Email, err.Error()))
		return
	}

	tokenStruct := AccessToken{
		AccessToken: accessToken,
	}

	stringForm, err := util.GetStringJson(tokenStruct)
	if err != nil{
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, stringForm)
}