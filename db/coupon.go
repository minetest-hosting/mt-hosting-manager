package db

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CouponRepository struct {
	g *gorm.DB
}

// coupon

func (r *CouponRepository) Insert(c *types.Coupon) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return r.g.Create(c).Error
}

func (r *CouponRepository) Update(n *types.Coupon) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *CouponRepository) GetByID(id string) (*types.Coupon, error) {
	return FindSingle[types.Coupon](r.g.Where(types.Coupon{ID: id}))
}

func (r *CouponRepository) GetByCode(code string) (*types.Coupon, error) {
	return FindSingle[types.Coupon](r.g.Where(types.Coupon{Code: code}))
}

func (r *CouponRepository) GetAll() ([]*types.Coupon, error) {
	return FindMulti[types.Coupon](r.g.Where(types.Coupon{}))
}

func (r *CouponRepository) GetRedeemedCoupons(id string) ([]*types.RedeemedCoupon, error) {
	return FindMulti[types.RedeemedCoupon](r.g.Where(types.RedeemedCoupon{CouponID: id}))
}

func (r *CouponRepository) Delete(id string) error {
	return r.g.Delete(types.Coupon{ID: id}).Error
}

func (r *CouponRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.Coupon{}).Error
}

// redeemed_coupon

func (r *CouponRepository) Redeem(coupon_id, user_id string) error {
	rc := &types.RedeemedCoupon{
		CouponID:  coupon_id,
		UserID:    user_id,
		Timestamp: time.Now().Unix(),
	}
	return r.g.Create(rc).Error
}

func (r *CouponRepository) IsRedeemed(coupon_id, user_id string) (bool, error) {
	entry, err := FindSingle[types.RedeemedCoupon](r.g.Where(types.RedeemedCoupon{CouponID: coupon_id, UserID: user_id}))
	return entry != nil, err
}
