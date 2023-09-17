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
//  @Success 200 {array} models.User
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := user.FindAll()

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("No users found"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, users)
}

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

	newUser, err := user.Create(userData)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Error creating user"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, newUser)
}

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
