// Root router of the application
package routes

import (
	"fmt"
	// "log"
	"net/http"
	// "os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	// "github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"server/authorization"
	"server/handlers"
	middlewareCustom "server/middleware"

	// "server/models"
	_ "server/docs"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = authorization.InitJWTAuth()
}

// Returns a router with all routes configured
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

	// Protected routes
	router.Group(func(r chi.Router) {
		router.Route("/api/v1/admin", func(r chi.Router) {
			//1. Verify token
			//2. Authenticate token
			//3. Populate roles from token into context
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator) //TODO: Add custom Authenticator to verify token against DB
			r.Use(middlewareCustom.RBACMiddleware)

			r.With(middlewareCustom.RBACMiddlewareProtectedRoute("admin")).Get("/", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				w.Write([]byte(fmt.Sprintf("Hello, %v you are authorized to view this.", claims["user_id"])))
			})
		})
	})

	return router
}
