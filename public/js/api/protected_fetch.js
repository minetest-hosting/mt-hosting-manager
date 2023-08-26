
export const store = Vue.reactive({
    title: "",
    message: "",
    url: "",
    status: 0
});

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
        store.title = `HTTP fetch error`;
        store.url = url;
        store.status = err.code;
        store.message = err.message;
    });
