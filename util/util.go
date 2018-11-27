package util

import (
"regexp"
"golang.org/x/crypto/bcrypt"
"encoding/json"
"fmt"
"net/http"
"time"
)

func IsValidEmail(email string) (bool) {
	if (len(email) <= 0) {
		return false
	}
	return regexp.MustCompile(`^([a-zA-Z0-9_\.\-\+])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$`).MatchString(email)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if (err != nil){
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Max(price1 float64, price2 float64)(float64){
	if price1 >= price2{
		return price1
	}
	return price2
}

func Min(price1 float64, price2 float64)(float64){
	if price1 >= price2{
		return price2
	}
	return price1
}

func GetNoDataSuccessResponse() string {
	return "{\"success\":1}"
}

func GetStringJson(v interface{})(string, error){
	jsonForm, err := json.Marshal(v)
	if (err != nil){
		return "", err
	}
	return string(jsonForm), nil
}

func CustomError(w http.ResponseWriter, error string, code int){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintf(w, error)
}

func CheckSameDate(time1 time.Time, time2 time.Time)(bool){
	if (time1.Year()==time2.Year())&&(time1.Month()==time2.Month())&&(time1.Day()==time2.Day()){
		return true;
	}
	return false
}

func IntAbsVal(x int) (int){
	if x < 0 {
		return -1 * x
	}
	return x
}
