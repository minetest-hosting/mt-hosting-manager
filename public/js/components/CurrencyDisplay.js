import { get_user_profile } from "../service/user.js";
import { get_rates } from "../service/currency.js";

export default {
    props: ["eurocents"],
    computed: {
        amount: function() {
            const rates = get_rates();
            const rate = parseFloat(rates[this.currency]);
            return Math.floor(this.eurocents / 100 * rate * 1000000) / 1000000;
        },
        currency: function() {
            const profile = get_user_profile();
            return profile.currency || "EUR";
        }
    },
    template: /*html*/`
        &euro; {{eurocents/100}} / {{amount}} {{currency}}
    `
};