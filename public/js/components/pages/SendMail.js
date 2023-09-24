import CardLayout from "../layouts/CardLayout.js";
import UserSearch from "../UserSearch.js";

import { send_mail } from "../../api/mail.js";

export default {
	components: {
		"card-layout": CardLayout,
		"user-search": UserSearch
	},
	methods: {
		send: function() {
			send_mail({
				subject: this.subject,
				content: this.content
			}, this.user.id);
		}
	},
	data: function() {
		return {
			user: null,
			subject: "",
			content: "",
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "envelope", name: "Send mail", link: "/sendmail"
			}]
		};
	},
	template: /*html*/`
	<card-layout title="Send mail" icon="envelope" :breadcrumb="breadcrumb">
		<div class="row">
			<div class="col-8">
				<user-search v-model="user"/>
			</div>
			<div class="col-4">
				<button class="btn btn-primary w-100" v-on:click="send" :disabled="!subject || !user || !content">
					Send
				</button>
			</div>
		</div>
		<hr>
		<div class="row">
			<div class="col-12">
				<input class="form-control" v-model="subject" placeholder="Subject"/>
			</div>
		</div>
		<div class="row">
			<div class="col-12">
				<textarea rows="12" class="form-control" v-model="content"></textarea>
			</div>
		</div>
	</card-layout>
	`
};
