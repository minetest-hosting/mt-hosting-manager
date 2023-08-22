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
			<h4>Minetest hosting</h4>
			<hr/>
			<img src="assets/minetest-hosting-600px.png" class="img img-rounded"/>
		</div>
	</card-layout>
	`
};
