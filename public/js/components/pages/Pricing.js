import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";
import NodeTypeSpec from "../NodeTypeSpec.js";

import { get_nodetypes } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay,
        "node-type-spec": NodeTypeSpec
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "money-bill", name: "Pricing", link: "/pricing"
            }]
		};
	},
    computed: {
        node_types: function(){
            return get_nodetypes().filter(nt => nt.state == "ACTIVE");
        }
    },
	template: /*html*/`
	<card-layout title="Pricing" icon="money-bill" :breadcrumb="breadcrumb">
        <h4>Available server-types</h4>
        <div class="card w-100" style="padding: 10px; margin-bottom: 10px;" v-for="nt in node_types">
            <div class="row">
                <div class="col-md-3">
                    <h4>
                        {{nt.name}}
                    </h4>
                    <h5>
                        {{nt.location_readable}}
                    </h5>
                    <node-type-spec :nodetype="nt"/>
                </div>
                <div class="col-md-4">
                    <ul>
                        <li>
                            Daily cost: <currency-display :eurocents="nt.daily_cost"/>
                        </li>
                        <li>
                            For 30 days: <currency-display :eurocents="nt.daily_cost * 30"/>
                        </li>
                    </ul>
                </div>
                <div class="col-md-5">
                    {{nt.description}}
                </div>
            </div>
        </div>
	</card-layout>
	`
};
