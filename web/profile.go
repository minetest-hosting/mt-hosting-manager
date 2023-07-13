package web

import (
	"mt-hosting-manager/types"
	"net/http"
)

type ProfileModel struct {
	User *types.User
}

func (ctx *Context) Profile(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	pm := &ProfileModel{}

	user, err := ctx.repos.UserRepo.GetByMail(c.Mail)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	pm.User = user

	ctx.tu.ExecuteTemplate(w, r, "profile.html", pm)
}
