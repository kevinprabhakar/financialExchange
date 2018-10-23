package handlers

import (
	"net/http"
	"financialExchange/util"
)

func AuthenticateUser(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := util.GetAccessTokenFromHeader(w, r)
		if err != nil{
			http.Error(w, err.Error(), 401)
			ServerLogger.ErrorMsg("Error while grabbing accessToken")
			return
		}

		_, err = util.VerifyAccessToken(accessToken)
		if err != nil{
			http.Error(w, err.Error(), 401)
			ServerLogger.ErrorMsg("Error while validating accessToken")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func NoAuthUser(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		h.ServeHTTP(w, r)
	})
}