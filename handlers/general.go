package handlers

import (
	"net/http"
	"fmt"
)

func HomeHealthCheck(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Server is running!")
}