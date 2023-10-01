import { get_user_profile } from "../service/user.js";
import { get_rates, get_currency_list } from "../service/currency.js";

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
            const rates = get_rates();
            const rate = parseFloat(rates[this.currency]);
            return Math.floor(this.eurocents / 100 * rate * 1000000) / 1000000;
        },
        currency: function() {
            const profile = get_user_profile();
            return profile.currency || "EUR";
        },
        decimals: function() {
            const c = get_currency_list().find(c => c.id == this.currency);
            return c ? c.decimals : 2;
        },
        warn: function() {
            if (!this.enable_warning) {
                return false;
            }
            const profile = get_user_profile();
            return profile.balance < 500;
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