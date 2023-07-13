package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"mt-hosting-manager/db"
	"mt-hosting-manager/tmpl"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web/oauth"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type OauthHandler struct {
	Impl     oauth.OauthImplementation
	UserRepo *db.UserRepository
	Config   *oauth.OAuthConfig
	BaseURL  string
	Type     types.UserType
	Tu       *tmpl.TemplateUtil
}

func SendJson(w http.ResponseWriter, o any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (h *OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		h.Tu.RenderError(w, r, 404, errors.New("no code found"))
		return
	}

	code := list[0]

	access_token, err := h.Impl.RequestAccessToken(code, h.BaseURL, h.Config)
	if err != nil {
		h.Tu.RenderError(w, r, 500, err)
		return
	}

	info, err := h.Impl.RequestUserInfo(access_token, h.Config)
	if err != nil {
		h.Tu.RenderError(w, r, 500, err)
		return
	}

	fmt.Printf("%v\n", info)

	if info.Email == "" {
		h.Tu.RenderError(w, r, 500, errors.New("empty email"))
		return
	}

	if info.ExternalID == "" {
		h.Tu.RenderError(w, r, 500, errors.New("empty external_id"))
		return
	}

	// check if there is already a user by that name
	user, err := h.UserRepo.GetByMail(info.Email)
	if err != nil {
		h.Tu.RenderError(w, r, 500, err)
		return
	}

	if user == nil {
		user = &types.User{
			Created:    time.Now().Unix(),
			Name:       info.Name,
			Mail:       info.Email,
			ExternalID: info.ExternalID,
			Type:       h.Type,
			Role:       types.UserRoleUser,
		}
		err = h.UserRepo.Insert(user)
		if err != nil {
			h.Tu.RenderError(w, r, 500, err)
			return
		}
		logrus.WithFields(logrus.Fields{
			"name":        user.Name,
			"type":        user.Type,
			"mail":        user.Mail,
			"external_id": user.ExternalID,
		}).Debug("created new user")
	}

	dur := time.Duration(24 * 180 * time.Hour)
	claims := &types.Claims{
		Mail: user.Mail,
		Role: user.Role,
	}

	token, err := h.Tu.CreateJWT(claims, dur)
	if err != nil {
		h.Tu.RenderError(w, r, 500, err)
		return
	}

	h.Tu.SetClaims(w, token, dur)

	target := h.BaseURL + "/profile"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
