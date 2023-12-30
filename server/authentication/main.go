// Description: This file contains the main authentication logic for the server.
package authentication

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/env"
	"server/helpers"
	"server/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

type UserAuth struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Scope    string `json:"scope,omitempty"`
}

var userModel models.User

// Generates a JWT token for the user
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var user UserAuth

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Error().Err(err).Msg("Error decoding user")
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(user)

	var validationErrors []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Err(err).Msg("Error validating user")

			validationErrors = append(validationErrors, err.Error())
		}

		if len(validationErrors) > 0 {
			var errorMessages string

			for _, errorMessage := range validationErrors {
				errorMessages += errorMessage + "\n"
			}

			helpers.ErrorJSON(w, errors.New(errorMessages), http.StatusBadRequest)
			return
		}
	}

	current_user, err := userModel.FindByEmail(user.Username)
	if err != nil {
		log.Error().Err(err).Msg("Error finding user")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//validate user credentials
	verified := helpers.ComparePasswords(current_user.Password, user.Password)
	if !verified {
		helpers.ErrorJSON(w, errors.New("Invalid Credentials Passed"), http.StatusBadRequest)
		return
	}

	// roles := []string{current_user.UserRole}

	if verified {
		//create token
		var token models.JWTClaims = models.JWTClaims{
			Email: user.Username,
			AppMetadata: models.AppMetadata{
				Authorization: models.Authorization{
					// Roles: roles,
				},
			},
			Subject:    user.Username,
			Audience:   "HOST", //TODO: Add audience from env
			Expiration: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:   time.Now().Unix(),
			//TODO: Generate JWTID from database
		}

		log.Info().Msgf("token: %v", token)

		//Create JWT token
		jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))

		jwtToken.Claims = token

		signedToken, err := jwtToken.SignedString([]byte(env.DefaultConfig.JWT_SECRET))

		if err != nil {
			log.Error().Err(err).Msg("Error signing token")
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		log.Info().Msgf("signedToken: %v", signedToken)

		response := struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{
			AccessToken:  signedToken,
			RefreshToken: "",
		}

		_ = helpers.WriteJSON(w, http.StatusOK, response)
		return
	}

	helpers.ErrorJSON(w, errors.New("Invalid Credentials Passed"), http.StatusBadRequest)
}
