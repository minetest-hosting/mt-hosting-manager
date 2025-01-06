package types

type Coupon struct {
	ID         string `json:"id" gorm:"primarykey;column:id"`
	Name       string `json:"name" gorm:"column:name"`
	Code       string `json:"code" gorm:"column:code"`
	ValidFrom  int64  `json:"valid_from" gorm:"column:valid_from"`
	ValidUntil int64  `json:"valid_until" gorm:"column:valid_until"`
	Value      int64  `json:"value" gorm:"column:value"`
	UseCount   int64  `json:"use_count" gorm:"column:use_count"`
	UseMax     int64  `json:"use_max" gorm:"column:use_max"`
}

func (*Coupon) TableName() string {
	return "coupon"
}

type RedeemedCoupon struct {
	CouponID  string `json:"coupon_id" gorm:"primarykey;column:coupon_id"`
	UserID    string `json:"user_id" gorm:"primarykey;column:user_id"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
}

func (*RedeemedCoupon) TableName() string {
	return "redeemed_coupon"
}
