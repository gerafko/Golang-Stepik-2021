package main

import (
	"html/template"
	"net/http"
	Handlers "redditclone/pkg/handlers"
	Posts "redditclone/pkg/items"
	"redditclone/pkg/middleware"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	Tmpl *template.Template
}

func main() {
	templates := template.Must(template.ParseGlob("./template/*.html"))

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	var sm *session.SessionJWT = session.NewSessionsJWT(session.Secret)

	var PostsRepo *Posts.PostsRepo = Posts.NewRepo()
	var userRepo *user.UserRepo = user.NewUserRepo()

	handlers := &Handlers.ItemsHandler{
		Tmpl:      templates,
		Logger:    logger,
		PostsRepo: PostsRepo,
	}

	userHandler := &Handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Index).Methods("GET")                  //- индекс
	r.HandleFunc("/api/login", userHandler.SignIn).Methods("POST")    //- логин
	r.HandleFunc("/api/register", userHandler.SignUp).Methods("POST") // - регистрация

	r.HandleFunc("/api/posts/", handlers.List).Methods("GET")                                                                         //- список всех постов
	r.Handle("/api/posts", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.AddPost))).Methods("POST")                         //- добавление поста - обратите внимание - есть с урлом, а есть с текстом
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", handlers.SpecList).Methods("GET")                                                      //- список постов конкретной категории
	r.HandleFunc("/api/post/{POST_ID}", handlers.PostDetail).Methods("GET")                                                           //- детали поста с комментами
	r.Handle("/api/post/{POST_ID}", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.AddComm))).Methods("POST")                //- добавление коммента
	r.Handle("/api/post/{POST_ID}/{COMMENT_ID}", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.DelComm))).Methods("DELETE") //- удаление коммента
	r.Handle("/api/post/{POST_ID}/upvote", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.Rate))).Methods("GET")             //- рейтинг поста вверх - ГЕТ был сделан автором оригинального фронта, я пока не добрался форкнуть и поправить.
	r.Handle("/api/post/{POST_ID}/downvote", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.Rate))).Methods("GET")           //- рейтинг поста вниз
	r.Handle("/api/post/{POST_ID}/unvote", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.Rate))).Methods("GET")             //- отмена ( удаление ) своего голоса в рейтинге
	r.Handle("/api/post/{POST_ID}", middleware.AuthMiddleware(sm, http.HandlerFunc(handlers.DeletePost))).Methods("DELETE")           //- удаление поста
	r.HandleFunc("/api/user/{USER_LOGIN}", handlers.UserPostList).Methods("GET")                                                      //- получение всех постов конкртеного пользователя
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./template")))

	addr := ":8080"
	http.ListenAndServe(addr, r)
}
