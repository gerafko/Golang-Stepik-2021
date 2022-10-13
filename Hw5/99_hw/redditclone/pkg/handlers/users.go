package Handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"go.uber.org/zap"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo *user.UserRepo
	Sessions *session.SessionJWT
}

func (repo *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := parsePayLoad(r)
	u, err := repo.UserRepo.Register(data.Login, data.Password)
	if err != nil {
		http.Error(w, "users.SignUp.repo.UserRepo.Register", http.StatusUnauthorized)
		return
	}
	s, err := repo.Sessions.Create(u)
	if err != nil {
		http.Error(w, "users.SignUp.repo.Sessions.Create", http.StatusUnauthorized)
		return
	}
	w.Write(s)
}

func (repo *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	data := parsePayLoad(r)
	u, err := repo.UserRepo.Authorize(data.Login, data.Password)
	if err != nil {
		http.Error(w, "users.SignIn.repo.UserRepo.Authorize", http.StatusUnauthorized)
		return
	}
	s, err := repo.Sessions.Create(u)
	if err != nil {
		http.Error(w, "users.SignIn.repo.Sessions.Create", http.StatusUnauthorized)
		return
	}
	w.Write(s)
}

func parsePayLoad(r *http.Request) user.User {
	var data user.User
	json.NewDecoder(r.Body).Decode(&data)
	return data
}
