package web

import (
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/components"
	"net/http"
)

type ProfileModel struct {
	User       *types.User
	Breadcrumb *components.Breadcrumb
}

func (ctx *Context) Profile(w http.ResponseWriter, r *http.Request, c *types.Claims) {

	user, err := ctx.repos.UserRepo.GetByMail(c.Mail)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &ProfileModel{
		User: user,
		Breadcrumb: &components.Breadcrumb{
			Entries: []*components.BreadcrumbEntry{
				components.HomeBreadcrumb,
				components.ProfileBreadcrumb,
			},
		},
	}

	ctx.tu.ExecuteTemplate(w, r, "profile.html", m)
}
