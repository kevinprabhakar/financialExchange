package main

import(
	"net/http"
	"github.com/gorilla/mux"
	"financialExchange/handlers"
	"fmt"
)

var port = ":3001"


func main(){
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHealthCheck).Methods("GET", "POST")

	r.HandleFunc("/api/entity", handlers.CreateEntity).Methods("POST")
	r.HandleFunc("/api/entity/login", handlers.SignInEntity).Methods("POST")

	r.HandleFunc("/api/customer", handlers.SignUpCustomer).Methods("POST")
	r.HandleFunc("/api/customer/login", handlers.SignInCustomer).Methods("POST")

	r.HandleFunc("/api/order", handlers.PlaceOrder).Methods("POST")

	fmt.Printf("Listening for requests on port %s\n", port)
	http.ListenAndServe(port, r)
}
