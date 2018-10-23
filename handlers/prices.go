package handlers

import (
	"net/http"
	"encoding/json"
	"financialExchange/model"
	"errors"
	"financialExchange/pricebook"
	"fmt"
)

func GetPrices(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.PriceRequest
	err := decoder.Decode(&params)

	if err != nil{
		http.Error(w, errors.New("InvalidPriceRequestParams").Error(), 400)
		return
	}

	prices, err := PriceController.GetSecurityChartForTimePeriod(pricebook.TimePeriod(params.TimePeriod), params.Security)
	if err != nil{
		http.Error(w, errors.New(fmt.Sprintf("Error Retreiving Prices for Security: %s", err.Error())).Error(), 500)
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	byteForm, err := json.Marshal(prices)
	if err != nil{
		http.Error(w, errors.New(fmt.Sprintf("Error Marshalling Price Data: %s", err.Error())).Error(), 500)
		return
	}

	fmt.Fprintf(w, string(byteForm))
}
