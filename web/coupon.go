package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RedeemStatus struct {
	Success bool `json:"success"`
}

func (a *Api) RedeemCoupon(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	code := vars["code"]

	coupon, err := a.repos.CouponRepo.GetByCode(code)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get coupon error: %v", err))
		return
	}
	if coupon == nil {
		// no such code
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}
	if coupon.UseCount >= coupon.UseMax {
		// max reached
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}
	now := time.Now().Unix()
	if now > coupon.ValidUntil || now < coupon.ValidFrom {
		// not valid yet or expired
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}

	already_redeemed, err := a.repos.CouponRepo.IsRedeemed(coupon.ID, c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get redeemed_coupon error: %v", err))
		return
	}
	if already_redeemed {
		// user already redeemed coupon
		Send(w, RedeemStatus{Success: false}, nil)
		return
	}

	// everything ok here!
	err = a.repos.UserRepo.AddBalance(c.UserID, coupon.Value)
	if err != nil {
		SendError(w, 500, fmt.Errorf("addbalance error: %v", err))
		return
	}

	// set redeemed
	err = a.repos.CouponRepo.Redeem(coupon.ID, c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("redeem error: %v", err))
		return
	}

	// increment use count
	coupon.UseCount++
	err = a.repos.CouponRepo.Update(coupon)
	if err != nil {
		SendError(w, 500, fmt.Errorf("update usecount error: %v", err))
		return
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:   types.AuditLogCouponRedeemed,
		UserID: c.UserID,
		Amount: &coupon.Value,
	})

	// all done
	Send(w, RedeemStatus{Success: true}, nil)
}

func (a *Api) CreateCoupon(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	coupon := &types.Coupon{}
	err := json.NewDecoder(r.Body).Decode(coupon)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	err = a.repos.CouponRepo.Insert(coupon)
	Send(w, coupon, err)
}
