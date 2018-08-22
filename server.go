package main

import(
	"net/http"
	"fmt"
)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Testing")
	})

	http.ListenAndServe(":8080", nil)
}
