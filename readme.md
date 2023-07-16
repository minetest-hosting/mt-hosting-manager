
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