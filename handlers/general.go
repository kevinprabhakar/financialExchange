package handlers

import (
	"net/http"
	"fmt"
	"financialExchange/util"
	"financialExchange/model"
	"encoding/json"
	"database/sql"
)

func HomeHealthCheck(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Server is running!")
}

func SearchPrefixes(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = util.VerifyAccessToken(accessToken)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var searchPrefix model.SearchParams
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&searchPrefix)

	searchResults := make([]model.SearchResult, 0)

	entities, err := TransactionController.Database.GetEntitiesWithPrefix(searchPrefix.Prefix)
	if err != nil{
		if err != sql.ErrNoRows{
			http.Error(w, err.Error(), 500)
			return
		}
	}
	securities, err := TransactionController.Database.GetSecuritiesWithPrefix(searchPrefix.Prefix)
	if err != nil{
		if err != sql.ErrNoRows{
			http.Error(w, err.Error(), 500)
			return
		}
	}
	for _, entity := range entities{
		security, err := TransactionController.Database.GetSecurityByID(entity.Security)
		if err != nil{
			http.Error(w, err.Error(), 500)
			return
		}
		searchResult := model.SearchResult{
			EntityName: entity.Name,
			Symbol: security.Symbol,
			SecurityID: entity.Security,
		}
		searchResults = append(searchResults, searchResult)
	}
	for _, security := range securities{
		entity, err := TransactionController.Database.GetEntityByID(security.Entity)
		if err != nil{
			http.Error(w, err.Error(), 500)
			return
		}
		searchResult := model.SearchResult{
			EntityName: entity.Name,
			Symbol: security.Symbol,
			SecurityID: security.Id,
		}
		searchResults = append(searchResults, searchResult)
	}

	byteForm, err := json.Marshal(searchResults)
	if err != nil{
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, string(byteForm))



}