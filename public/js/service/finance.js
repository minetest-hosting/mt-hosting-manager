import { get_balance } from "./user.js";

export function get_refund_amount(tx) {
    if (tx.amount_refunded > 0) {
        return 0;
    }
    const balance = get_balance();
    if (balance <= 0) {
        return 0;
    }
    return Math.min(tx.amount, balance);
}