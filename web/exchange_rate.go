package web

import (
	"net/http"
)

func (a *Api) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	list, err := a.repos.ExchangeRateRepo.GetAll()
	Send(w, list, err)
}
