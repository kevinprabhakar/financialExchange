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

	orderCreateErr := OrderController.CreateOrder(orderCreateParams)

	if orderCreateErr != nil{
		ServerLogger.ErrorMsg(orderCreateErr.Error())
		http.Error(w, orderCreateErr.Error(), 500)
		return
	}

	fmt.Fprint(w, util.GetNoDataSuccessResponse())

}
