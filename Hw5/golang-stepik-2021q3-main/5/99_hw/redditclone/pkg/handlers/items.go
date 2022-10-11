package Handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	Posts "redditclone/pkg/items"
	Auth "redditclone/pkg/middleware"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ItemsHandler struct {
	Tmpl      *template.Template
	PostsRepo *Posts.PostsRepo
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *ItemsHandler) SpecList(w http.ResponseWriter, r *http.Request) {
	spec, exist := getFromURL(r, "CATEGORY_NAME")
	if !exist {
		message := fmt.Sprintf("Error SpecList.CATEGORY_NAME is NULL: %s", spec)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	a, err := h.PostsRepo.GetSpec(spec)
	if err != nil {
		message := fmt.Sprintf("Error SpecList.GetSpec(): %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(a)
	if err != nil {
		message := fmt.Sprintf("Error SpecList.Marshal(NewRepo()): %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	a, err := h.PostsRepo.GetAll()
	if err != nil {
		message := fmt.Sprintf("Error List.GetAll(): %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(a)
	if err != nil {
		message := fmt.Sprintf("Error List.NewRepo(): %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (h *ItemsHandler) PostDetail(w http.ResponseWriter, r *http.Request) {
	id, ok := getFromURL(r, "POST_ID")
	if ok != true {
		message := fmt.Sprintf("Error PostDetail.POST_ID is NULL: id=%s", id)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	a, err := h.PostsRepo.GetByPostId(id)
	if err != nil {
		message := fmt.Sprintf("Error PostDetail.GetByID(): %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(a)
	w.Write(b)
}

func (h *ItemsHandler) UserPostList(w http.ResponseWriter, r *http.Request) {
	user_login, exist := getFromURL(r, "USER_LOGIN")
	if !exist {
		message := fmt.Sprintf("Error UserPostList.USER_LOGIN is NULL: USER_LOGIN=%s", user_login)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	a, err := h.PostsRepo.GetByLogin(user_login)
	if err != nil {
		message := fmt.Sprintf("Error UserPostList.GetByLogin(): %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(a)
	w.Write(b)
}

func (h *ItemsHandler) Rate(w http.ResponseWriter, r *http.Request) {
	data := parsePayLoadItem(r)
	a := r.Context().Value(Auth.SessionKey).(*session.Session)
	if a == nil {
		message := fmt.Sprintf("ErrNoAuth")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	id, ok := getFromURL(r, "POST_ID")
	if !ok {
		message := fmt.Sprintf("Error AddComm.POST_ID is NULL")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	data.Author.ID = a.UserId
	data.Author.Username = a.UserName
	post, err := h.PostsRepo.Rate(data, id)
	if err != nil {
		message := fmt.Sprintf("Rate Error 1: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(post)
	if err != nil {
		message := fmt.Sprintf("Rate Error 2: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
}

func getFromURL(r *http.Request, teg string) (string, bool) {
	vars := mux.Vars(r)
	spec, exist := vars[teg]
	return spec, exist
}

func (h *ItemsHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	data := parsePayLoadItem(r)
	a := r.Context().Value(Auth.SessionKey).(*session.Session)
	if a == nil {
		message := fmt.Sprintf("ErrNoAuth")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	data.Author.ID = a.UserId
	data.Author.Username = a.UserName
	post, err := h.PostsRepo.AddPost(data)
	if err != nil {
		message := fmt.Sprintf("AddPost Error Adding: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(post)
	if err != nil {
		message := fmt.Sprintf("AddPost Error Marshaling: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
}

func (h *ItemsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, ok := getFromURL(r, "POST_ID")
	if !ok {
		message := fmt.Sprintf("Error DeletePost.POST_ID is NULL")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	post, err := h.PostsRepo.DelPost(id)
	if err != nil {
		message := fmt.Sprintf("AddPost Error Adding: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(post)
	if err != nil {
		message := fmt.Sprintf("AddPost Error Marshaling: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
}

func (h *ItemsHandler) AddComm(w http.ResponseWriter, r *http.Request) {
	data := parsePayLoadComment(r)
	a := r.Context().Value(Auth.SessionKey).(*session.Session)
	if a == nil {
		message := fmt.Sprintf("ErrNoAuth")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	id, ok := getFromURL(r, "POST_ID")
	if !ok {
		message := fmt.Sprintf("Error AddComm.POST_ID is NULL")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	post, err := h.PostsRepo.AddComm(data, id, a.UserId, a.UserName)
	if err != nil {
		message := fmt.Sprintf("AddComm Error Adding: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(post)
	if err != nil {
		message := fmt.Sprintf("AddComm Error Marshaling: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
}

func (h *ItemsHandler) DelComm(w http.ResponseWriter, r *http.Request) {
	//data := parsePayLoadItem(r)
	id, ok := getFromURL(r, "POST_ID")
	comm_id, ok := getFromURL(r, "COMMENT_ID")
	if !ok {
		message := fmt.Sprintf("Error DelComm.POST_ID is NULL")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	a := r.Context().Value(Auth.SessionKey).(*session.Session)
	if a == nil {
		message := fmt.Sprintf("ErrNoAuth")
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	post, err := h.PostsRepo.DelComm(id, comm_id)
	if err != nil {
		message := fmt.Sprintf("DelComm Error Adding: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(post)
	if err != nil {
		message := fmt.Sprintf("DelComm Error Marshaling: %s", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
}

func parsePayLoadItem(r *http.Request) Posts.Post {
	var data Posts.Post
	json.NewDecoder(r.Body).Decode(&data)
	return data
}

func parsePayLoadComment(r *http.Request) Posts.IncomingComment {
	var data Posts.IncomingComment
	json.NewDecoder(r.Body).Decode(&data)
	return data
}
