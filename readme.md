
## TODO

* [ ] Login options
  * [x] github
  * [ ] discord
  * [ ] mesehub
  * [ ] gitlab
* [ ] Host creation
  * [ ] Hetzner
  * [ ] Contabo
* [ ] Host provisioning (ssh)
* [ ] Instance creation
* [ ] Instance setup (mtui)
* [ ] Backups?

# Dev

```sh
# set all users as admin
sudo sqlite3 mt-hosting.sqlite "update user set role = 'ADMIN'"
```

# Environment variables

* `LOGLEVEL` "debug" / "info"
* `ENABLE_WORKER`

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

* `WALLEE_USERID`
* `WALLEE_SPACEID`
* `WALLEE_KEY`

* `NTFY_URL`
* `NTFY_TOPIC`
* `NTFY_USERNAME`
* `NTFY_PASSWORD`