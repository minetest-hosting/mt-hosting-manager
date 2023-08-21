// app bootstrap script
// hostname == "localhost": debug bundles
// else: prod bundles

let scripts = [];

if (location.hostname == "localhost") {
    // dev
    let script = document.createElement("script");
    script.src = "node_modules/vue/dist/vue.global.js";
    scripts.push(script);
    
    script = document.createElement("script");
    script.src = "node_modules/vue-router/dist/vue-router.global.js";
    scripts.push(script);

} else {
    // prod
    let script = document.createElement("script");
    script.src = "node_modules/vue/dist/vue.global.prod.js";
    scripts.push(script);
    
    script = document.createElement("script");
    script.src = "node_modules/vue-router/dist/vue-router.global.prod.js";
    scripts.push(script);
}

// local code
let script = document.createElement("script");
script.src = "js/bundle.js";
script.onerror = function() {
    // dev/module fallback if bundle is not found
    let script = document.createElement("script");
    script.src = "js/main.js";
    script.type = "module";
    document.body.appendChild(script);
};
scripts.push(script);

function load_script() {
    if (scripts.length == 0) {
        return;
    }

    let s = scripts.shift();
    s.onload = load_script;
    document.body.appendChild(s);
}

load_script();