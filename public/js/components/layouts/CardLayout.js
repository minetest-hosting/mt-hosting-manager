import Breadcrumb from "../Breadcrumb.js";

export default {
    props: ["title", "icon", "breadcrumb"],
    components: {
        "bread-crumb": Breadcrumb
    },
	template: /*html*/`
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8">
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
