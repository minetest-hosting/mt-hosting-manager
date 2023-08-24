import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-server"></i> Node details
		</template>
	</card-layout>
	`
};
