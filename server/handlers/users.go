package handlers

import (
	"encoding/json"
	"errors"

	// "log"
	// "io/ioutil"
	"net/http"

	"server/helpers"
	"server/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	// "server/types"
	// "github.com/go-chi/chi/v5"
)

var user models.User

// Get All Users
//
//	 @Summary      Get all Users
//	 @Description  Get all Users
//	 @Tags         users
//	 @Accept       json
//	 @Produce      json
//	 @Router       /api/v1/users [get]
//	 @Success 200 {array} models.User
//	 @Failure 500 {object} string
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := user.FindAll()

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("No users found"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, users)
}

// Create User
//
//	 @Summary      Create User
//	 @Description  Create User
//	 @Tags         users
//	 @Accept       json
//	 @Produce      json
//	 @Router       /api/v1/users [post]
//	 @Success 200 {object} models.User
//	 @Failure 500 {object} string
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userData models.User

	// log.Println("Body: ", r.Body)
	helpers.MessageLogs.InfoLog.Println("Body: ", r.Body)

	err := json.NewDecoder(r.Body).Decode(&userData)
	// body, err := ioutil.ReadAll(r.Body)
	// err = json.Unmarshal(body, &userData)

	helpers.MessageLogs.InfoLog.Println("userData", userData)

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

	newUser, err := user.Create(userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Error creating user"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, newUser)
}

// Find User By Email
//
//	 @Summary      Find User By Email
//	 @Description  Find User By Email
//	 @Tags         users
//	 @Accept       json
//	 @Produce      json
//	 @Router       /api/v1/users/{email} [get]
//	 @Param email path string true "Email"
//	 @Success 200 {object} models.User
//	 @Failure 500 {object} string
func FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

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
//	 @Summary      Update User By Email
//	 @Description  Update User By Email
//	 @Tags         users
//	 @Accept       json
//	 @Produce      json
//	 @Router       /api/v1/users [put]
//	 @Success 200 {object} models.User
//	 @Failure 500 {object} string
func UpdateUserByEmail(w http.ResponseWriter, r *http.Request) {
	var userData models.User

	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Invalid JSON"), http.StatusBadRequest)
		return
	}

	err = user.UpdateByEmail(userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("No user found"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, user)
}
