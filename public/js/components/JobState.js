export default {
    props: ["state"],
	template: /*html*/`
	<span class="badge bg-success" v-if="state == 'CREATED'">
		<i class="fa fa-pause"></i>
		Created
	</span>
	<span class="badge bg-success" v-if="state == 'RUNNING'">
		<i class="fa fa-play"></i>
		Running
	</span>
	<span class="badge bg-success" v-if="state == 'DONE_SUCCESS'">
		<i class="fa fa-check"></i>
		Success
	</span>
	<span class="badge bg-success" v-if="state == 'DONE_FAILURE'">
		<i class="fa fa-times"></i>
		Failed
	</span>
	`
};
