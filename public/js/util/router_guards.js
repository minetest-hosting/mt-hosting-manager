import login_store from '../store/login.js';

const LoginPath = { path: "/login" };

export default function(router) {
    router.beforeEach((to) => {
        if (to.meta.loggedIn && !login_store.loggedIn) {
            return LoginPath;
        }

        if (to.meta.requiredRole) {
            if (!login_store.loggedIn) {
                // quick login check
                return LoginPath;
            }

            if (login_store.claims.role != to.meta.requiredRole){
                // check required role
                return LoginPath;
            }
        }
    });   
}