package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggo/http-swagger/v2"

	"server/handlers"
	// "server/models"
	_ "server/docs"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is up and running"))
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"), //The url pointing to API definition
	))

	router.Group(func(r chi.Router) {
		router.Route("/api/v1/users", func(r chi.Router) {
			r.Get("/", handlers.GetAllUsers)
			r.Post("/", handlers.CreateUser)
			r.Get("/{email}", handlers.FindUserByEmail)
			r.Put("/", handlers.UpdateUserByEmail)
		})
	})

	// router.Get("/api/v1/users", handlers.GetAllUsers)
	// router.Post("/api/v1/users", handlers.CreateUser)

	return router
}
