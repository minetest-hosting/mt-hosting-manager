import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-play"></i> Jobs
		</template>
	</card-layout>
	`
};
