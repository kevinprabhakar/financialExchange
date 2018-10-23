package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"financialExchange/util"
	"financialExchange/model"
	"strconv"
	gosql "database/sql"
	"errors"
	"github.com/gorilla/mux"
	"time"
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

func GetCurrEntity(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	uid, err := util.VerifyAccessToken(accessToken)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	intform, err := strconv.ParseInt(uid, 10, 64)
	if err != nil{
		http.Error(w, errors.New("InvalidEntityID").Error(), 400)
		return
	}

	entity, err := EntityController.Database.GetEntityByID(intform)
	if err != nil{
		http.Error(w, errors.New("Couldn't get customer").Error(), 500)
		return
	}

	EntityController.Logger.Debug(fmt.Sprintf("Retrieved entity profile for %s", entity.Email))


	byteForm, err := json.Marshal(entity)
	if err != nil{
		http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
		return
	}

	fmt.Fprintf(w, string(byteForm))
}

func CreateSecurity(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	uid, err := util.VerifyAccessToken(accessToken)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	intform, err := strconv.ParseInt(uid, 10, 64)
	if err != nil{
		http.Error(w, errors.New("InvalidEntityID").Error(), 400)
		return
	}

	//Verify security does not already exist for this security
	security, err := EntityController.Database.GetSecurityByEntityID(intform)
	if err != nil{
		//If error is not a no rows error, return http error
		if err != gosql.ErrNoRows{
			http.Error(w, errors.New(fmt.Sprintf("Error Retreiving Security From DB: %s", err.Error())).Error(), 500)
			return
		}
	}else{
		http.Error(w, errors.New(fmt.Sprintf("Security Already Exists for Entity: %s", security.Symbol)).Error(), 400)
		return
	}

	vars := mux.Vars(r)
	symbol, ok := vars["symbol"]
	if !ok{
		http.Error(w, errors.New("Invalid Security Name").Error(), 400)
		return
	}

	newSecurity := model.Security{
		Entity: intform,
		Symbol: symbol,
		Created:time.Now().Unix(),
	}

	_, err = EntityController.Database.InsertSecurityIntoTable(newSecurity)
	if err != nil{
		ServerLogger.ErrorMsg(fmt.Sprintf("Error while creating Security %s: %s", security.Symbol, err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	ServerLogger.Debug(fmt.Sprintf("Created Security %s: %s", security.Symbol, err.Error()))


	fmt.Fprintf(w, util.GetNoDataSuccessResponse())
}