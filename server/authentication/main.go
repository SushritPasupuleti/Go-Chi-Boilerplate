// Description: This file contains the main authentication logic for the server.
package authentication

import (
	"errors"
	"github.com/go-chi/oauth"
	"net/http"
)

type UserVerifier struct {
}

// VerifyUser is a method that verifies a user's credentials
func (*UserVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	//TODO: implement user validation from database
	if len(username) > 0 && len(password) > 0 {
		return nil
	}

	return errors.New("wrong user")
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (*UserVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	//TODO: implement client validation from database
	if len(clientID) > 0 && len(clientSecret) > 0 {
		return nil
	}

	return errors.New("wrong client")
}

// ValidateCode validates token ID
func (*UserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (*UserVerifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	claims := make(map[string]string)
	//TODO: Add claims from database
	return claims, nil
}

// AddProperties provides additional information to the token response
func (*UserVerifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	//TODO: Add properties from database
	return props, nil
}

// ValidateTokenID validates token ID
func (*UserVerifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*UserVerifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
