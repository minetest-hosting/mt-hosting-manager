import { protected_fetch } from "./protected_fetch.js";

export const redeem_coupon = code => protected_fetch(`api/coupon/redeem/${code}`, {
    method: "POST"
});

export const create_coupon = coupon => protected_fetch(`api/coupon`, {
    method: "POST",
    body: JSON.stringify(coupon)
});
