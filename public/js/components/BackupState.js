export default {
    props: ["state"],
	template: /*html*/`
	<span class="badge bg-info" v-if="state == 'CREATED'">
		<i class="fa fa-pause"></i>
		Created
	</span>
	<span class="badge bg-primary" v-if="state == 'PROGRESS'">
		<i class="fa fa-spinner fa-spin"></i>
		In progress
	</span>
	<span class="badge bg-success" v-if="state == 'COMPLETE'">
		<i class="fa fa-check"></i>
		Complete
	</span>
	<span class="badge bg-danger" v-if="state == 'ERROR'">
		<i class="fa fa-times"></i>
		Error
	</span>
	`
};
