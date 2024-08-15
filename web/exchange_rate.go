package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	list, err := a.repos.ExchangeRateRepo.GetAll()
	Send(w, list, err)
}

func (a *Api) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	list, err := a.repos.ExchangeRateRepo.GetByCurrency(vars["currency"])
	Send(w, list, err)
}

func (a *Api) CreateExchangeRate(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	rate := &types.ExchangeRate{}
	err := json.NewDecoder(r.Body).Decode(rate)
	if err != nil {
		SendError(w, 500, fmt.Errorf("error parsing json: %v", err))
		return
	}

	err = a.repos.ExchangeRateRepo.Insert(rate)
	if err != nil {
		SendError(w, 500, fmt.Errorf("error inserting rate: %v", err))
		return
	}
	Send(w, rate, nil)
}

func (a *Api) UpdateExchangeRate(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	rate := &types.ExchangeRate{}
	err := json.NewDecoder(r.Body).Decode(rate)
	if err != nil {
		SendError(w, 500, fmt.Errorf("error parsing json: %v", err))
		return
	}

	err = a.repos.ExchangeRateRepo.Update(rate)
	if err != nil {
		SendError(w, 500, fmt.Errorf("error inserting rate: %v", err))
		return
	}
	Send(w, rate, nil)
}

func (a *Api) DeleteExchangeRate(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	err := a.repos.ExchangeRateRepo.DeleteByCurrency(vars["currency"])
	Send(w, true, err)
}
