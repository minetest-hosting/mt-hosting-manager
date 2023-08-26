import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-home"></i> Template
		</template>
	</card-layout>
	`
};
