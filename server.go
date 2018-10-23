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
	r.Handle("/api/security", handlers.AuthenticateUser(handlers.GetSecurity)).Methods("GET")

	//todo: Fill out the rest of these CRUD functions
	//orders
	r.Handle("/api/customer/orders", handlers.AuthenticateUser(handlers.GetUserOrders)).Methods("GET")
	r.Handle("/api/customer/orders/{orderID}", handlers.AuthenticateUser(handlers.GetUserOrders)).Methods("GET")

	//customer
	r.Handle("/api/customer", handlers.AuthenticateUser(handlers.GetUser)).Methods("GET")

	//portfolio
	r.Handle("/api/customer/portfolio", handlers.AuthenticateUser(handlers.GetCurrUserPortfolio)).Methods("GET")

	//ownedShares
	r.Handle("/api/customer/ownedShares", handlers.AuthenticateUser(handlers.GetCurrUserOwnedShares)).Methods("GET")
	r.Handle("/api/customer/ownedShares/{security}", handlers.AuthenticateUser(handlers.GetCurrUserOwnedSharesForSecurity)).Methods("GET")

	//entity
	r.Handle("/api/entity", handlers.AuthenticateUser(handlers.GetCurrEntity)).Methods("GET")

	//prices
	//in the future authenticate for prices, but for testing purposes, we will keep it unauthenticated
	r.Handle("/api/prices", handlers.NoAuthUser(handlers.GetPrices)).Methods("POST")


	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/src")))

	fmt.Printf("Listening for requests on port %s\n", port)
	http.ListenAndServe(port, r)
}