import App from './app.js';
import routes from './routes.js';
import { check_login } from './service/login.js';
import { fetch_exchange_rates } from './service/exchange_rate.js';
import { fetch_info } from './service/info.js';
import router_guards from './util/router_guards.js';
import events, { EVENT_STARTUP } from './events.js';
import "./service/user.js";

function start(){
	// create router instance
	const router = VueRouter.createRouter({
		history: VueRouter.createWebHashHistory(),
		routes: routes
	});

	// set up router guards
	router_guards(router);

	// trigger startup event
	events.emit(EVENT_STARTUP);

	// start vue
	const app = Vue.createApp(App);
	app.use(router);
	app.mount("#app");
}

fetch_info()
.then(() => Promise.all([check_login(), fetch_exchange_rates()]))
.then(() => start())
.catch(e => console.error(e));
