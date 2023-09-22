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
	// "server/redis"

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
		log.Error().Err(err).Msg("Error getting users")
		helpers.ErrorJSON(w, errors.New("No users found"), http.StatusInternalServerError)
		return
	}

	//save to cache
	middleware.SaveToCache(r, users)

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

	// log.Info().Msgf("Body: %t", r.Body)

	err := json.NewDecoder(r.Body).Decode(&userData)
	// body, err := ioutil.ReadAll(r.Body)
	// err = json.Unmarshal(body, &userData)

	log.Info().Msgf("userData: %v", userData)

	if err != nil {
		log.Error().Err(err).Msg("Error decoding JSON")
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(userData)

	var validationErrors []string

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Err(err).Msg("Error validating user")

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
		log.Error().Err(err).Msg("Error creating user")

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

	validate := validator.New()

	err := validate.Var(email, "required,email")

	if err != nil {
		log.Error().Err(err).Msg("Error validating email")
		helpers.ErrorJSON(w, errors.New("Invalid email"), http.StatusBadRequest)
		return
	}

	user, err := user.FindByEmail(email)

	if err != nil {
		log.Error().Err(err).Msg("Error finding user")
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
		log.Error().Err(err).Msg("Error decoding JSON")
		helpers.ErrorJSON(w, errors.New("Invalid JSON"), http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(userData)

	var validationErrors []string

	if err != nil {

		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Err(err).Msg("Error validating user")

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
		log.Error().Err(err).Msg("Error updating user")
		helpers.ErrorJSON(w, errors.New("No user found"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, user)
}
