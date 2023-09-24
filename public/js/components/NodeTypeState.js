export default {
    props: ["state"],
	template: /*html*/`
	<span class="badge bg-success" v-if="state == 'ACTIVE'">
		Active
	</span>
	<span class="badge bg-info" v-if="state == 'INACTIVE'">
		Inactive
	</span>
	<span class="badge bg-warning" v-if="state == 'DEPRECATED'">
		Deprecated
	</span>
	`
};
