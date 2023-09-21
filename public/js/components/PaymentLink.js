
export default {
    props: ["id"],
    template: /*html*/`
    <router-link :to="'/finance/detail/' + id">
        <i class="fa fa-money-bill"></i>
        {{id}}
    </router-link>
	`
};
