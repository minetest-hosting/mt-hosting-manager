import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { get_nodetypes } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay
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
                <div class="col-3">
                    <h4>
                        {{nt.name}}
                    </h4>
                    <h5>
                        {{nt.location_readable}}
                    </h5>
                </div>
                <div class="col-2">
                    <ul>
                        <li>
                            <i class="fa-solid fa-microchip"></i> CPU: {{nt.cpu_count}}
                            <span class="badge bg-success" v-if="nt.dedicated">
                                Dedicated
                            </span>
                        </li>
                        <li>
                            <i class="fa-solid fa-hard-drive"></i> Disk: {{nt.disk_gb}} GB
                        </li>
                        <li>
                            <i class="fa-solid fa-clipboard"></i> RAM: {{nt.ram_gb}} GB
                        </li>
                    </ul>
                </div>
                <div class="col-2">
                    <ul>
                        <li>
                            Daily cost: <currency-display :eurocents="nt.daily_cost"/>
                        </li>
                        <li>
                            For 30 days: <currency-display :eurocents="nt.daily_cost * 30"/>
                        </li>
                    </ul>
                </div>
                <div class="col-5">
                    {{nt.description}}
                </div>
            </div>
        </div>
	</card-layout>
	`
};
