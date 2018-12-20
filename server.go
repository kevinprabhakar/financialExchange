package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"financialExchange/handlers"
	"fmt"
)

var port = ":3001"


func main() {
	r := mux.NewRouter()

	//r.HandleFunc("/", handlers.HomeHealthCheck).Methods("GET", "POST")

	r.HandleFunc("/api/entity", handlers.CreateEntity).Methods("POST")
	r.HandleFunc("/api/entity/login", handlers.SignInEntity).Methods("POST")
	r.HandleFunc("/api/entity/security/{symbol}", handlers.CreateSecurity).Methods("POST")
	r.HandleFunc("/api/entity/ipo", handlers.IPO).Methods("POST")

	r.HandleFunc("/api/customer", handlers.SignUpCustomer).Methods("POST")
	r.HandleFunc("/api/customer/login", handlers.SignInCustomer).Methods("POST")

	r.Handle("/api/order", handlers.AuthenticateUser(handlers.PlaceOrder)).Methods("POST")
	r.Handle("/api/security/{symbol}", handlers.AuthenticateUser(handlers.GetSecurity)).Methods("GET")
	r.Handle("/api/securities/mostdailytraded", handlers.AuthenticateUser(handlers.GetMostDailyTraded)).Methods("GET")

	//todo: Fill out the rest of these CRUD functions
	//orders
	r.Handle("/api/customer/orders", handlers.AuthenticateUser(handlers.GetUserOrders)).Methods("GET")
	r.Handle("/api/customer/orders/{orderID}", handlers.AuthenticateUser(handlers.GetUserOrders)).Methods("GET")
	r.Handle("/api/customer/giveMoney", handlers.AuthenticateUser(handlers.GiveUserMoney)).Methods("POST")


	//customer
	r.Handle("/api/customer", handlers.AuthenticateUser(handlers.GetUser)).Methods("GET")

	//portfolio
	r.Handle("/api/customer/portfolio", handlers.AuthenticateUser(handlers.GetCurrUserPortfolio)).Methods("GET")

	//ownedShares
	r.Handle("/api/customer/ownedShares", handlers.AuthenticateUser(handlers.GetCurrUserOwnedShares)).Methods("GET")
	r.Handle("/api/customer/ownedShares/{security}", handlers.AuthenticateUser(handlers.GetCurrUserOwnedSharesForSecurity)).Methods("GET")

	//entity
	r.Handle("/api/entity", handlers.AuthenticateUser(handlers.GetCurrEntity)).Methods("GET")
	r.Handle("/api/entity/{symbol}", handlers.AuthenticateUser(handlers.GetEntityBySymbol)).Methods("GET")

	//search autocompletion
	r.Handle("/api/search", handlers.AuthenticateUser(handlers.SearchPrefixes)).Methods("POST")


	//prices
	//in the future authenticate for prices, but for testing purposes, we will keep it unauthenticated
	r.Handle("/api/prices", handlers.NoAuthUser(handlers.GetPrices)).Methods("POST")
	r.Handle("/api/currprice/{symbol}", handlers.NoAuthUser(handlers.GetCurrPrice)).Methods("GET")


	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/src")))

	fmt.Printf("Listening for requests on port %s\n", port)
	http.ListenAndServe(port, r)
}