import Breadcrumb from "../Breadcrumb.js";

export default {
    props: ["title", "icon", "breadcrumb", "fullwidth"],
    components: {
        "bread-crumb": Breadcrumb
    },
    computed: {
        center_col_classes: function() {
            if (this.fullwidth) {
                return { "col-12": true };
            } else {
                return { "col-8": true };
            }
        }
    },
	template: /*html*/`
    <div class="row">
        <div class="col-2" v-if="!fullwidth"></div>
        <div v-bind:class="center_col_classes">
            <bread-crumb :items="breadcrumb" v-if="breadcrumb"/>
            <div class="card">
                <div class="card-header">
                    <i v-bind:class="{fa:true, ['fa-'+icon]:true}"></i> {{title}}
                </div>
                <div class="card-body">
                    <slot></slot>
                </div>
            </div>
        </div>
    </div>
	`
};
