//Custom Middlware for RBAC and Authorization
package middleware

import (
	"context"
	"log"
	"net/http"

	"server/authorization"

	"server/helpers"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = authorization.InitJWTAuth()
}

// Extracts the claims from the JWT token and adds them to the context
func RBACMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		ctx := context.WithValue(r.Context(), "claims", claims)

		scope := claims["app_metadata"].(map[string]interface{})["authorization"] //.(map[string]interface{})["roles"]

		log.Printf("RBACMiddleware: scope=%v\n", scope)

		var scopeArray []string

		//extract roles from scope
		for key, value := range scope.(map[string]interface{}) {
			// log.Printf("RBACMiddleware: key=%v, value=%v type=%T\n", key, value, value)

			valueArray := value.([]interface{})
			// log.Printf("RBACMiddleware: valueArray=%v\n", valueArray)

			if key == "roles" {
				for _, v := range valueArray {
					// log.Printf("RBACMiddleware: v=%v\n", v)
					scopeArray = append(scopeArray, v.(string))
				}
			}
		}

		log.Printf("RBACMiddleware: scopeArray=%v\n", scopeArray)

		ctx = context.WithValue(ctx, "scope", scopeArray)

		//TODO: Add other claims to context

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Checks if the user has the required scope to access the route
func RBACMiddlewareProtectedRoute(scopeRequired string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// log.Println("RBACMiddleware: before request")
			// log.Printf("Context: %v\n", ctx)
			// log.Printf("RBACMiddleware: user=%v\n", ctx.Value("user"))

			scope := r.Context().Value("scope").([]string)
			// scope := r.Context().Value("scope").([]interface{})

			log.Printf("RBACMiddlewareProtectedRoute: scope=%v\n", scope)
			log.Printf("RBACMiddlewareProtectedRoute: scopeRequired=%v\n", scopeRequired)

			if !helpers.Contains(scope, scopeRequired) {
				log.Printf("RBACMiddlewareProtectedRoute: scopeRequired=%v not found in scope=%v\n", scopeRequired, scope)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
