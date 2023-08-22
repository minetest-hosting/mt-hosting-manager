import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-home"></i> Start
		</template>
		<div class="text-center">
			<h3>
				Minetest hosting
			</h3>
		</div>
	</card-layout>
	`
};
