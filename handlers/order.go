package handlers

import (
	"net/http"
	"encoding/json"
	"financialExchange/model"
	"fmt"
	"financialExchange/util"
	"strconv"
	"errors"
)

func PlaceOrder(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var orderCreateParams model.OrderCreateParams
	err := decoder.Decode(&orderCreateParams)

	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	action := "sell"

	if orderCreateParams.InvestorAction == 0{
		action = "buy"
	}

	ServerLogger.Debug(fmt.Sprintf("Placing %s order of %d shares at $%.2f for Customer ID=%d", action,orderCreateParams.NumShares, orderCreateParams.CostPerShare, orderCreateParams.UserID))

	orderInsertID, orderCreateErr := OrderController.CreateOrder(orderCreateParams)

	if orderCreateErr != nil{
		ServerLogger.ErrorMsg(orderCreateErr.Error())
		http.Error(w, orderCreateErr.Error(), 500)
		return
	}

	orderIdStruct := model.OrderId{
		Id: orderInsertID,
	}

	stringForm, err := util.GetStringJson(orderIdStruct)
	if err != nil{
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, stringForm)
}

func IPO(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var IPOParams model.IPOParams
	err := decoder.Decode(&IPOParams)

	if err != nil{
		ServerLogger.ErrorMsg(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

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

	_, err = OrderController.IPO(IPOParams, intform)
	if err != nil{
		http.Error(w, errors.New("IPO for Entity failed").Error(), 400)
		return
	}

	ServerLogger.Debug(fmt.Sprintf("Successful IPO for %s", uid))

	fmt.Fprintf(w, util.GetNoDataSuccessResponse())
}