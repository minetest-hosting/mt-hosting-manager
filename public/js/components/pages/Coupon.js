import CardLayout from "../layouts/CardLayout.js";

import { get_coupons, create_coupon } from "../../api/coupon.js";

import format_time from "../../util/format_time.js";

function new_coupon() {
    return {
        name: "",
        code: "",
        value: 0,
        valid_from: new Date(),
        valid_until: new Date(Date.now() + (1000 * 3600 * 24 * 30)),
        use_count: 0,
        use_max: 10
    };
}

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "ticket", name: "Coupon", link: "/coupon"
            }],
            new_coupon: new_coupon(),
            list: []
		};
	},
    created: function() {
        this.update();
    },
    methods: {
        format_time,
        update: async function() {
            this.list = await get_coupons();
        },
        create: async function() {
            await create_coupon(Object.assign({}, this.new_coupon, {
                valid_from: Math.floor(+this.new_coupon.valid_from/1000),
                valid_until: Math.floor(+this.new_coupon.valid_until/1000)
            }));
            await this.update();
            this.new_coupon = new_coupon();
        }
    },
	template: /*html*/`
	<card-layout title="Coupons" icon="ticket" :breadcrumb="breadcrumb" :fullwidth="true">
        <table class="table table-condensed table-striped">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Code</th>
                    <th>Value</th>
                    <th>Valid from</th>
                    <th>Valid until</th>
                    <th>Use count</th>
                    <th>Max use</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="coupon in list">
                    <td>{{coupon.name}}</td>
                    <td>{{coupon.code}}</td>
                    <td>{{coupon.value}}</td>
                    <td>{{format_time(coupon.valid_from)}}</td>
                    <td>{{format_time(coupon.valid_until)}}</td>
                    <td>{{coupon.use_count}}</td>
                    <td>{{coupon.use_max}}</td>
                    <td>
                    </td>
                </tr>
                <tr>
                    <td>
                        <input class="form-control" type="text" v-model="new_coupon.name" placeholder="Name"/>
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="new_coupon.code" placeholder="Code"/>
                    </td>
                    <td>
                        <input class="form-control" type="number" min="0" step="1" v-model="new_coupon.value" placeholder="Value in eurocents"/>
                    </td>
                    <td>
                        <vue-datepicker v-model="new_coupon.valid_from"/>
                    </td>
                    <td>
                        <vue-datepicker v-model="new_coupon.valid_until"/>
                    </td>
                    <td>
                    </td>
                    <td>
                        <input class="form-control" type="number" min="1" step="1" v-model="new_coupon.use_max" placeholder="Max use count"/>
                    </td>
                    <td>
                        <button class="btn btn-primary" v-on:click="create">
                            <i class="fa fa-plus"></i>
                            Add
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
	</card-layout>
	`
};
