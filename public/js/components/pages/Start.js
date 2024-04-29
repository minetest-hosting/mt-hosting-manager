import CardLayout from "../layouts/CardLayout.js";
import Breadcrumb from "../Breadcrumb.js";

export default {
	components: {
		"card-layout": CardLayout,
		"bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [{icon: "home", name: "Home", link: "/"}]
		};
	},
	template: /*html*/`
	<card-layout title="Home" icon="home" :breadcrumb="breadcrumb">
		<div class="text-center">
			<h4>Minetest hosting</h4>
			<img src="assets/minetest-hosting-80px.png" class="img img-rounded"/>
			<hr/>
			<div class="alert alert-warning">
				<i class="fa-solid fa-triangle-exclamation"></i>
				<b>Warning:</b> This project is still in an experimental stage, feel free to try things but don't expect stability at this point
			</div>
			<router-link class="btn btn-secondary" to="/pricing">
				<i class="fa fa-money-bill"></i> Pricing
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="mailto:hosting@minetest.ch">
				<i class="fa fa-envelope"></i> Contact
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://discord.gg/Xj62fUbQkn" target="new">
				<i class="fa-brands fa-discord"></i>
				Join the discord server
			</a>
			&nbsp;
			<router-link class="btn btn-outline-secondary" to="/privacy-policy">
				<i class="fa fa-section"></i> Privacy policy
			</router-link>
			&nbsp;
			<router-link class="btn btn-outline-secondary" to="/terms-conditions">
				<i class="fa fa-section"></i> Terms and conditions
			</router-link>
			&nbsp;
			<a class="btn btn-outline-secondary" href="https://github.com/minetest-hosting/mt-hosting-manager" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
			&nbsp;
			<a class="btn btn-outline-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Powered by mtui
			</a>
		</div>
		<hr/>
		<div class="text-center">
			<h4>Features</h4>
		</div>
		<div class="row">
			<div class="col-6">
				<h5>Easy server provisioning</h5>
				<ul>
					<li>Create servers with just a few clicks</li>
					<li>IPv6 support by default</li>
					<li>Secured with TLS (the management interface)</li>
					<li>Free *.minetest.ch domain included or bring own</li>
					<li>Pay with crypto or ordinary currency (not enabled yet)</li>
				</ul>
			</div>
			<div class="col-6">
				<img src="assets/features_create_server.png" class="img img-rounded"/>
			</div>
		</div>
		<hr>
		<div class="row">
			<div class="col-6">
				<h5>Server management web-interface</h5>
				<ul>
					<li>Use the wizard to set up your server or change every setting you want</li>
					<li>Easy mod management with ContentDB and Git support</li>
					<li>Player- and Ban-Management out of the box</li>
					<li>Integrations for SkinsDB and ingame Mail available</li>
					<li>Execute chat-commands or inject lua-code through the web-interface</li>
					<li>Chat with your players</li>
				</ul>
			</div>
			<div class="col-6">
				<img src="assets/feature_mod_management.png" class="img img-rounded"/>
			</div>
		</div>
		<hr>
		<div class="row">
			<div class="col-6">
				<h5>Wide choice of servers and locations</h5>
				<ul>
					<li>Available locations: Germany 🇩🇪 and United States 🇺🇸</li>
					<li>Select your server according to demand, see <router-link to="/pricing">Pricing</router-link></li>
					<li>Run multiple minetest-engines on a single server</li>
				</ul>
			</div>
			<div class="col-6">
				<img src="assets/features_node_types.png" class="img img-rounded"/>
			</div>
		</div>
	</card-layout>
	`
};
