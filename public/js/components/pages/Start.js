import store from "../../store/info.js";

export default {
	data: () => store,
	template: /*html*/`
	<div>
		<div class="text-center">
			<h3>
				Minetest hosting
			</h3>
		</div>
		<a :href="'https://github.com/login/oauth/authorize?client_id=' + GithubClientID + '&scope=user:email'" class="btn btn-secondary">
            <i class="fab fa-github"></i>
            Login with Github
        </a>
	</div>
	`
};
