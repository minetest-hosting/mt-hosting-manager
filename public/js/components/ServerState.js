export default {
    props: ["state"],
	template: /*html*/`
	<span class="badge bg-success" v-if="state == 'CREATED'">
		<i class="fa fa-pause"></i>
		Created
	</span>
	<span class="badge bg-success" v-if="state == 'PROVISIONING'">
		<i class="fa fa-spinner fa-spin"></i>
		Provisioning
	</span>
	<span class="badge bg-success" v-if="state == 'RUNNING'">
		<i class="fa fa-play"></i>
		Running
	</span>
	<span class="badge bg-success" v-if="state == 'REMOVING'">
		<i class="fa fa-times"></i>
		Removing
	</span>
	`
};
