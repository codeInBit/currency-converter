package rates

import (
	"net/http"
	"os"
	"strconv"

	"github.com/codeinbit/currency-converter/utilities"
	"github.com/gorilla/mux"
)

var baseurl = "https://v6.exchangerate-api.com/v6/"

var SupportedCurrencies = []string{"EUR", "USD", "GBP", "AUD", "JPY", "HKD"}

func isSupported(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func UnitRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fromCurrency := vars["from"]
	toCurrency := vars["to"]

	url := baseurl + os.Getenv("API_KEY") + "/pair/" + fromCurrency + "/" + toCurrency

	if !isSupported(SupportedCurrencies, fromCurrency) {
		utilities.ERROR(w, http.StatusBadRequest, "Invalid currency symbol, '"+fromCurrency+"'")
		return
	}

	if !isSupported(SupportedCurrencies, toCurrency) {
		utilities.ERROR(w, http.StatusBadRequest, "Invalid currency symbol, '"+toCurrency+"'")
		return
	}

	conversionRate, responseDuration := fetchRate(url, fromCurrency+"-"+toCurrency)

	utilities.JSON(w, http.StatusOK, conversionRate, true, responseDuration)
}

func RateOnAmount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fromCurrency := vars["from"]
	toCurrency := vars["to"]
	amount := vars["amount"]
	amountInFloat, _ := strconv.ParseFloat(amount, 64)

	url := baseurl + os.Getenv("API_KEY") + "/pair/" + fromCurrency + "/" + toCurrency

	if !isSupported(SupportedCurrencies, fromCurrency) {
		utilities.ERROR(w, http.StatusBadRequest, "Invalid currency symbol, '"+fromCurrency+"'")
		return
	}

	if !isSupported(SupportedCurrencies, toCurrency) {
		utilities.ERROR(w, http.StatusBadRequest, "Invalid currency symbol, '"+toCurrency+"'")
		return
	}

	conversionRate, responseDuration := fetchRate(url, fromCurrency+"-"+toCurrency)
	conversionResult := conversionRate * amountInFloat

	utilities.JSON(w, http.StatusOK, conversionResult, true, responseDuration)
}
