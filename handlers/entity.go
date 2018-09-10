package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"financialExchange/util"
	"financialExchange/model"
)

func CreateEntity(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.CreateEntityParams
	err := decoder.Decode(&params)

	if err != nil{
		http.Error(w, err.Error(), 400)
		return
	}
	ServerLogger.Debug(fmt.Sprintf("Creating Entity %s", params.Name))
	err = EntityController.CreateEntity(params)

	if err != nil{
		ServerLogger.ErrorMsg(fmt.Sprintf("Error while creating Entity %s: %s", params.Name, err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, util.GetNoDataSuccessResponse())
}

func SignInEntity(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.SignInEntityParams
	err := decoder.Decode(&params)

	if err != nil{
		http.Error(w, err.Error(), 400)
		return
	}

	ServerLogger.Debug(fmt.Sprintf("Signing In %s", params.Email))

	accessToken, err := EntityController.SignIn(params)

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