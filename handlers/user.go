package handlers

import (
	"net/http"
	"encoding/json"
	"financialExchange/util"
	"fmt"
	"financialExchange/model"
	"github.com/gorilla/mux"
	"errors"
	"strconv"
)

func SignUpCustomer(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.CustomerSignUpParams
	err := decoder.Decode(&params)

	if err != nil{
		ServerLogger.ErrorMsg(err.Error())

		http.Error(w, err.Error(), 400)
		return
	}
	ServerLogger.Debug(fmt.Sprintf("Signing Up %s", params.Email))
	err = CustomerController.SignUp(params)

	if err != nil{
		http.Error(w, err.Error(), 500)
		ServerLogger.ErrorMsg(fmt.Sprintf("Error while signing up %s: %s", params.Email, err.Error()))
		return
	}

	SignInParams := model.CustomerSignInParams{
		Email: params.Email,
		Password: params.Password,
	}

	accessToken, err := CustomerController.SignIn(SignInParams)

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

func SignInCustomer(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var params model.CustomerSignInParams
	err := decoder.Decode(&params)

	if err != nil{
		http.Error(w, err.Error(), 400)
		return
	}

	ServerLogger.Debug(fmt.Sprintf("Signing In %s", params.Email))

	accessToken, err := CustomerController.SignIn(params)

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

func GetUser(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	customer, err := CustomerController.GetCurrUser(accessToken)
	if err != nil{
		http.Error(w, errors.New("Couldn't get customer").Error(), 500)
		return
	}

	CustomerController.Logger.Debug(fmt.Sprintf("Retrieved customer profile for %s", customer.Email))


	byteForm, err := json.Marshal(customer)
	if err != nil{
		http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
		return
	}

	fmt.Fprintf(w, string(byteForm))
}

func GetCurrUserPortfolio(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	portfolio, err := CustomerController.GetCurrUserPortfolio(accessToken)
	if err != nil{
		http.Error(w, errors.New("Couldn't get customer portfolio").Error(), 500)
		ServerLogger.ErrorMsg(err.Error())
		return
	}

	CustomerController.Logger.Debug(fmt.Sprintf("Retrieved customer portfolio for %d", portfolio.Customer))


	byteForm, err := json.Marshal(portfolio)
	if err != nil{
		http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
		return
	}

	fmt.Fprintf(w, string(byteForm))
}

func GetUserOrders(w http.ResponseWriter, r *http.Request){
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

	uidIntForm, err := strconv.ParseInt(uid, 10, 64)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	orderID, ok := mux.Vars(r)["orderID"]
	if !ok{
		orders, err := CustomerController.GetOrdersForUser(accessToken)

		if err != nil{
			http.Error(w, errors.New("Could not get orders for user").Error(), 500)
			return
		}
		byteForm, err := json.Marshal(*orders)
		if err != nil{
			http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
			return
		}

		fmt.Fprintf(w, string(byteForm))
	}else{
		int64Version, err := strconv.ParseInt(orderID, 10, 64)
		if err != nil{
			http.Error(w, errors.New("Couldn't parse orderID").Error(), 400)
			return
		}
		order, err := CustomerController.Database.GetOrderById(int64Version)
		if err != nil{
			http.Error(w, errors.New("Error getting order by orderID").Error(), 500)
			return
		}
		if order.Investor != uidIntForm{
			http.Error(w, errors.New("UserAttemptingToAccessNonUserOrder").Error(), 400)
			return
		}
		byteForm, err := json.Marshal(*order)
		if err != nil{
			if err != nil{
				http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
				return
			}
		}
		fmt.Fprintf(w, string(byteForm))
	}
}

func GetCurrUserOwnedShares(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	uid, err := util.VerifyAccessToken(accessToken)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	int64Version, err := strconv.ParseInt(uid, 10, 64)
	if err != nil{
		http.Error(w, errors.New("Couldn't parse user id").Error(), 400)
		return
	}

	ownedShares, err := CustomerController.Database.GetAllOwnedSharesForUserID(int64Version)
	if err != nil{
		http.Error(w, errors.New("Couldn't grab user owned shares").Error(), 400)
		return
	}

	byteForm, err := json.Marshal(ownedShares)
	if err != nil{
		http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
		return
	}

	fmt.Fprintf(w, string(byteForm))
}

func GetCurrUserOwnedSharesForSecurity(w http.ResponseWriter, r *http.Request){
	accessToken, err := util.GetAccessTokenFromHeader(w, r)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	uid, err := util.VerifyAccessToken(accessToken)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	securityId, ok := vars["security"]
	if !ok{
		http.Error(w, errors.New("No security Name Provided").Error(), http.StatusUnauthorized)
		return
	}

	uidint64, err := strconv.ParseInt(uid, 10, 64)
	if err != nil{
		http.Error(w, errors.New("Couldn't parse user id").Error(), 400)
		return
	}

	securityIdInt64, err := strconv.ParseInt(securityId, 10, 64)
	if err != nil{
		http.Error(w, errors.New("Couldn't parse security id").Error(), 400)
		return
	}

	ownedShares, err := CustomerController.Database.GetOwnedShareForUserForSecurity(uidint64, securityIdInt64)
	if err != nil{
		http.Error(w, errors.New("Couldn't grab user owned shares").Error(), 400)
		return
	}

	byteForm, err := json.Marshal(ownedShares)
	if err != nil{
		http.Error(w, errors.New("Error Marshalling JSON").Error(), 400)
		return
	}

	fmt.Fprintf(w, string(byteForm))
}