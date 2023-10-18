import { get_user_profile } from "../service/user.js";
import { get_rate } from "../service/exchange_rate.js";

function get_currency_sign(currency) {
    switch (currency) {
    case "EUR":
        return "€";
    case "USD":
        return "$";
    default:
        return currency;
    }
}

function format_currency(amount, currency, decimals) {
    const sign = get_currency_sign(currency);
    return `${sign} ${amount.toFixed(decimals)}`;
}

export default {
    props: ["eurocents", "enable_warning"],
    methods: {
        format_currency
    },
    computed: {
        amount: function() {
            const r = get_rate(this.currency);
            const rate = parseFloat(r.rate);
            return Math.floor(this.eurocents / 100 * rate * 1000000) / 1000000;
        },
        currency: function() {
            const profile = get_user_profile();
            return profile.currency || "EUR";
        },
        decimals: function() {
            return get_rate(this.currency).digits;
        },
        warn: function() {
            if (!this.enable_warning) {
                return false;
            }
            const profile = get_user_profile();
            return profile.balance < profile.warn_balance;
        }
    },
    template: /*html*/`
    <span class="badge" v-bind:class="{'bg-secondary': !warn, 'bg-warning': warn}">
        {{format_currency(eurocents/100, 'EUR', 2)}}
        <span v-if="currency != 'EUR'">
            / {{format_currency(amount, currency, decimals)}}
        </span>
    </span>
    `
};