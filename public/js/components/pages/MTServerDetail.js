import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-list"></i> Server details
		</template>
	</card-layout>
	`
};
