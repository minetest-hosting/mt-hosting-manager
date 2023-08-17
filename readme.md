
# Minetest hosting orchestrator

![](https://github.com/minetest-hosting/mt-hosting-manager/workflows/test/badge.svg)
![](https://github.com/minetest-hosting/mt-hosting-manager/workflows/build/badge.svg)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minetest-hosting/mt-hosting-manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/minetest-hosting/mt-hosting-manager)](https://goreportcard.com/report/github.com/minetest-hosting/mt-hosting-manager)
[![Coverage Status](https://coveralls.io/repos/github/minetest-hosting/mt-hosting-manager/badge.svg)](https://coveralls.io/github/minetest-hosting/mt-hosting-manager)


State: **WIP**

## Roadmap

* [ ] Login options
  * [x] github (MVP)
  * [ ] discord
  * [ ] mesehub
  * [ ] gitlab
* [ ] Payment options
  * [x] Wallee (MVP)
  * [ ] Stripe
  * [ ] Coinbase crypto
* [ ] Host creation
  * [x] Hetzner (MVP)
  * [ ] Contabo
* [x] Host provisioning (MVP)
* [ ] Instance setup
* [ ] Backups?

# Dev

```sh
# setup
docker-compose up
# set all users as admin
sudo sqlite3 mt-hosting.sqlite "update user set role = 'ADMIN'"
```

# Environment variables

* `LOGLEVEL` "debug" / "info"
* `ENABLE_WORKER`
* `STAGE` "prod" / "dev"
* `MOCK_ORCHESTRATION` mock orchestration part if "true"

* `CSRF_KEY`
* `JWT_KEY`
* `BASEURL`
* `WEBDEV`
* `COOKIE_DOMAIN`
* `COOKIE_PATH`
* `COOKIE_SECURE`

* `GITHUB_CLIENTID`
* `GITHUB_SECRET`

* `ADMIN_USER_MAIL` mail of the user that gets the admin role on register
* `DISABLE_SIGNUP`

* `WALLEE_USERID`
* `WALLEE_SPACEID`
* `WALLEE_KEY`

* `NTFY_URL`
* `NTFY_TOPIC`
* `NTFY_USERNAME`
* `NTFY_PASSWORD`

* `SSH_KEY`

# License

* Code: `MIT`
