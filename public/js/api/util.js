import error_toast_store from "../store/error_toast.js";

export const protected_fetch = (url, opts) => fetch(url, opts)
    .then(r => {
        if (r.status == 200) {
            return r.json();
         } else {
            return r.json().then(err => {
                return Promise.reject(err);
            });
         }
    })
    .catch(err => {
        error_toast_store.title = `HTTP fetch error`;
        error_toast_store.url = url;
        error_toast_store.status = err.code;
        error_toast_store.message = err.message;
    });