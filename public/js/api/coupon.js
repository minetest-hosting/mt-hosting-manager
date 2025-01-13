import { protected_fetch } from "./protected_fetch.js";

// admin api

export const get_coupon = id => protected_fetch(`api/coupon/${id}`);

export const get_coupon_users = id => protected_fetch(`api/coupon/${id}/users`);

export const get_coupons = () => protected_fetch(`api/coupon`);

export const update_coupon = coupon => protected_fetch(`api/coupon/${coupon.id}`, {
    method: "POST",
    body: JSON.stringify(coupon)
});

export const create_coupon = coupon => protected_fetch(`api/coupon`, {
    method: "POST",
    body: JSON.stringify(coupon)
});

// user api

export const redeem_coupon = code => protected_fetch(`api/coupon/redeem/${code}`, {
    method: "POST"
});

