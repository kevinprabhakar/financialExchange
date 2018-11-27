package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"strconv"
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