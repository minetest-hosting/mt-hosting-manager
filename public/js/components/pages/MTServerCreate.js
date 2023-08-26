import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-list"></i>
			<i class="fa fa-plus"></i>
			Create server
		</template>
	</card-layout>
	`
};
