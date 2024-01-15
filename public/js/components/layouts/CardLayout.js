import Breadcrumb from "../Breadcrumb.js";

export default {
    props: ["title", "icon", "breadcrumb", "fullwidth", "flex"],
    components: {
        "bread-crumb": Breadcrumb
    },
    computed: {
        center_col_classes: function() {
            if (this.fullwidth) {
                return { "col-md-12": true };
            } else {
                return { "col-md-8": true };
            }
        },
        card_styles: function() {
            if (this.flex) {
                return {
                    display: "flex",
                    "flex-wrap": "wrap"
                };
            } else {
                return {};
            }
        }
    },
	template: /*html*/`
    <div class="row">
        <div class="col-md-2" v-if="!fullwidth"></div>
        <div v-bind:class="center_col_classes">
            <bread-crumb :items="breadcrumb" v-if="breadcrumb"/>
            <div class="card">
                <div class="card-header">
                    <i v-bind:class="{fa:true, ['fa-'+icon]:true}"></i> {{title}}
                </div>
                <div class="card-body" v-bind:style="card_styles">
                    <slot></slot>
                </div>
            </div>
        </div>
    </div>
	`
};
