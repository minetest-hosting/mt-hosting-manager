import { redeem_coupon } from "../api/coupon.js";
import { fetch_profile } from "../service/user.js";

export default {
    data: function() {
        return {
            code: "",
            invalid_code: false,
            success: false
        };
    },
    methods: {
        redeem: async function() {
            this.invalid_code = false;
            this.success = false;

            const result = await redeem_coupon(this.code);
            if (result.success) {
                this.success = true;
                await fetch_profile();
            } else {
                this.invalid_code = true;
            }
        }
    },
    template: /*html*/ `
        <label>Enter coupon code</label>
        <div class="input-group">
            <input class="form-control" type="text" v-model="code" v-bind:class="{'is-invalid':invalid_code}"/>
            <button class="btn btn-outline-primary" v-on:click="redeem" placeholder="Code" :disabled="!code">
                <i class="fa fa-check" style="color: green;" v-if="success"></i>
                Redeem
            </button>
            <div class="invalid-feedback" v-if="invalid_code">
                Code invalid or used up
            </div>
        </div>
    `
};