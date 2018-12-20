package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func GetSecurity(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	symbol, ok := vars["symbol"]
	if !ok{
		http.Error(w, "No Symbol Provided", 400)
		return
	}


	if _, err := strconv.Atoi(symbol); err == nil {
		int64Form, err := strconv.ParseInt(symbol, 10, 64)
		if err != nil{
			http.Error(w, "Invalid Security ID", 400)
			ServerLogger.ErrorMsg(err.Error())
			return
		}
		security, err := OrderController.Database.GetSecurityByID(int64Form)
		if err != nil{
			http.Error(w, "No Security With This Name", 400)
			ServerLogger.ErrorMsg(err.Error())
			return
		}
		jsonForm, err := json.Marshal(security)
		if err != nil{
			http.Error(w, "Error Marshalling JSON", 500)
			return
		}

		fmt.Fprintf(w, string(jsonForm))
	}else{
		security, err := OrderController.Database.GetSecurityBySymbol(symbol)
		ServerLogger.Debug(symbol)
		if err != nil{
			http.Error(w, "No Security With This Name", 400)
			ServerLogger.ErrorMsg(err.Error())
			return
		}
		jsonForm, err := json.Marshal(security)
		if err != nil{
			http.Error(w, "Error Marshalling JSON", 500)
			return
		}

		fmt.Fprintf(w, string(jsonForm))
	}




}

func GetMostDailyTraded(w http.ResponseWriter, r *http.Request){
	startTime := time.Now().Add(-1 * 24 * time.Hour)

	mostTradedSecurities, err := OrderController.GetMostOrderedSecuritiesOverTimeframe(startTime, 10)
	if err != nil{
		http.Error(w, "Couldnt get most traded securities", 500)
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	jsonForm, err := json.Marshal(mostTradedSecurities)
	if err != nil{
		http.Error(w, "Error Marshalling JSON", 500)
		ServerLogger.ErrorMsg(err.Error())

		return
	}

	fmt.Fprintf(w, string(jsonForm))
}