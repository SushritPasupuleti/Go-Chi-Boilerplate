package handlers

import (
	// "context"
	"encoding/json"
	"errors"
	"fmt"
	// "time"

	// "go/format"

	// "log"
	// "io/ioutil"
	"net/http"

	"server/helpers"
	"server/middleware"
	"server/models"
	"server/redis"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	// "server/types"
	// "github.com/go-chi/chi/v5"
)

var user models.User

// Get All Users
//
//	@Summary      Get all Users
//	@Description  Get all Users
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Router       /api/v1/users [get]
//	@Success 200 {array} models.User
//	@Failure 500 {object} string
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := user.FindAll()

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("No users found"), http.StatusInternalServerError)
		return
	}

	r_ := r

	middleware.SaveToCache(r_, users)
	routeKey, err := middleware.PrepareRouteKey(r_)
	if err != nil {
		log.Error().Msgf("Error preparing route key: %v", err)
	}
	cacheKey, err := middleware.PrepareCacheKey(r.Body, routeKey)
	if err != nil {
		log.Error().Msgf("Error preparing cache key: %v", err)
	}

	cachedResponseRAW, err := redis.GetCache(cacheKey)
	if err != nil {
		log.Error().Msgf("Error getting cache: %v", err)
	}

	log.Info().Msgf("cachedResponseRAW: %v", cachedResponseRAW)

	cachedResponse, err := middleware.CachedResponseToJSON(cacheKey)
	if err != nil {
		log.Error().Msgf("Error getting cache: %v", err)
	}

	if err == nil {
		log.Info().Msgf("Sending CachedResponse %v", cachedResponse)
		helpers.WriteJSON(w, http.StatusOK, cachedResponse)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, users)
}

// Create User
//
//	@Summary      Create User
//	@Description  Create User
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Router       /api/v1/users [post]
//	@Success 200 {object} models.User
//	@Failure 500 {object} string
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userData models.User

	// log.Println("Body: ", r.Body)
	// helpers.MessageLogs.InfoLog.Println("Body: ", r.Body)
	// log.Info().Msgf("Route: %s", r.URL.Path)
	// log.Info().Msgf("Method: %s", r.Method)
	// log.Info().Msgf("Body: %t", r.Body)

	// key, errKey := middleware.PrepareCacheKey(r.Body, r.URL.Path, r.Method)
	//
	// if errKey != nil {
	// 	log.Error().Msgf("Error preparing cache key: %v", errKey)
	// }

	// log.Info().Msgf("key: %v", key)

	err := json.NewDecoder(r.Body).Decode(&userData)
	// body, err := ioutil.ReadAll(r.Body)
	// err = json.Unmarshal(body, &userData)

	log.Info().Msgf("userData: %v", userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(userData)

	var validationErrors []string

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			helpers.MessageLogs.ErrorLog.Println("Error: ", err)

			validationErrors = append(validationErrors, err.Error())
		}

		if len(validationErrors) > 0 {

			var errorMessages string

			for _, err := range validationErrors {
				errorMessages += err + "\n"
			}

			helpers.ErrorJSON(w, errors.New(errorMessages), http.StatusBadRequest)
			return
		}
	}

	newUser, err := user.Create(userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)

		errorMessage := fmt.Sprintf("Error creating user - Reason: %v", err)

		// helpers.ErrorJSON(w, errors.New("Error creating user \nReason: "), http.StatusInternalServerError)
		helpers.ErrorJSON(w, errors.New(errorMessage), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, newUser)
}

// Find User By Email
//
//	@Summary      Find User By Email
//	@Description  Find User By Email
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Router       /api/v1/users/{email} [get]
//	@Param email path string true "Email"
//	@Success 200 {object} models.User
//	@Failure 500 {object} string
func FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	// queryParams := r.URL.Query()

	// log.Info().Msgf("urlParams: %v", r.URL.RawQuery)
	// log.Info().Msgf("urlParams: %v", r.URL.Path)

	// log.Info().Msgf("queryParams: %v", queryParams)

	// ctx, cancel := context.WithTimeout(context.Background(), 5)
	// set, errx := redisClient.SetNX(ctx, "key", "value", 10*time.Second).Result()
	//
	// log.Info().Msgf("set: %v", set)
	// log.Info().Msgf("errx: %v", errx)

	validate := validator.New()

	err := validate.Var(email, "required,email")

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Invalid email"), http.StatusBadRequest)
		return
	}

	user, err := user.FindByEmail(email)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("No user found"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, user)
}

// Update User By Email
//
//	@Summary      Update User By Email
//	@Description  Update User By Email
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Router       /api/v1/users [put]
//	@Success 200 {object} models.User
//	@Failure 500 {object} string
func UpdateUserByEmail(w http.ResponseWriter, r *http.Request) {
	var userData models.User

	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Invalid JSON"), http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(userData)

	var validationErrors []string

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			helpers.MessageLogs.ErrorLog.Println("Error: ", err)

			validationErrors = append(validationErrors, err.Error())
		}

		if len(validationErrors) > 0 {

			var errorMessages string

			for _, err := range validationErrors {
				errorMessages += err + "\n"
			}

			helpers.ErrorJSON(w, errors.New(errorMessages), http.StatusBadRequest)
			return
		}
	}

	err = user.UpdateByEmail(userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("No user found"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, user)
}
