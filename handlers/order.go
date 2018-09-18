package handlers

import (
	"net/http"
	"encoding/json"
	"financialExchange/model"
	"fmt"
	"financialExchange/util"
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
