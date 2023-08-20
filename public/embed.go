package public

import (
	"embed"
)

//go:embed assets/*
//go:embed js/* index.html
//go:embed node_modules/bootswatch/dist/flatly/bootstrap.min.css
//go:embed node_modules/vue/dist/vue.global.prod.js
//go:embed node_modules/vue-router/dist/vue-router.global.prod.js
//go:embed node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed node_modules/@fortawesome/fontawesome-free/webfonts/*
var Webapp embed.FS
