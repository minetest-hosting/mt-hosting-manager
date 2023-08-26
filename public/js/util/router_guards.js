import { has_role, is_logged_in } from '../service/login.js';

const LoginPath = { path: "/login" };

export default function(router) {
    router.beforeEach((to) => {
        if (to.meta.loggedIn && !is_logged_in()) {
            return LoginPath;
        }

        if (to.meta.requiredRole && !has_role(to.meta.requiredRole)){
            // check required role
            return LoginPath;
        }
    });   
}