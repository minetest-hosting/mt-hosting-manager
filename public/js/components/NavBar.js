import { logout } from '../service/login.js';
import login_store from '../store/login.js';

export default {
	data: function() {
		return {
			login: login_store
		};
	},
	methods: {
		logout: function(){
			logout().then(() => this.$router.push("/login"));
		}
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Minetest hosting</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0" v-if="login.loggedIn">
					<li class="nav-item">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Home
						</router-link>
					</li>
				</ul>
				<div class="d-flex">
					<stats-display class="navbar-text" style="padding-right: 10px;"/>
					<div class="btn-group">
						<button class="btn btn-outline-secondary" v-if="login.claims">
							<router-link to="/profile">
								<i class="fas fa-user"></i>
								<span>
									Logged in as <b>{{login.claims.username}}</b>
								</span>
							</router-link>
						</button>
						<button class="btn btn-secondary" v-on:click="logout" v-if="login.loggedIn">
							<i class="fa-solid fa-right-from-bracket"></i>
							Logout
						</button>
					</div>
				<div>
			</div>
		</nav>
	`
};
