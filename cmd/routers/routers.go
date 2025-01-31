package routers

import (
	"database/sql"
	"net/http"

	"forum/backend/handlers"
	"forum/backend/handlers/auth"
	"forum/backend/handlers/posts"
	"forum/middleware"
)

func Router(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.IndexPage(w, r, db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(w, r, db)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(w, r, db)
	})

	http.HandleFunc("/sign-in", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginPage(w, r, db)
	})

	http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterPage(w, r, db)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.Logout(w, r, db)
	})

	http.HandleFunc("/addPost", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			posts.AddPost(w, r, db)
		}), db))

	http.HandleFunc("/addComment", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			posts.AddComments(w, r, db)
		}), db))

	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		posts.FilterByCat(w, r, db)
	})

	http.HandleFunc("/addlikeComment", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			posts.AddlikeComment(w, r, db)
		}), db))

	http.HandleFunc("/myPost", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			posts.FilterByUser(w, r, db)
		}), db))
	http.HandleFunc("/addlike", middleware.Authorization(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			posts.AddLikePost(w, r, db)
		}), db))
	http.HandleFunc("/frontend/static/", handlers.StaticHandler)
}
