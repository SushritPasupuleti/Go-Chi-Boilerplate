// Middleware and utility functions for caching requests
package middleware

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"server/redis"

	"github.com/rs/zerolog/log"
)

// Prepare route key for caching
func PrepareRouteKey(r *http.Request) (string, error) {

	return r.Method + "." + r.URL.Path + "." + r.URL.RawQuery, nil
}

// Returns a base64 encoded string of the payload with the route and method prepended
func PrepareCacheKey(payload interface{}, routeKey string) (string, error) {

	//format the payload into a string
	stringPayload := ""

	byteData, err := io.ReadAll(payload.(io.Reader))

	if err != nil {
		log.Fatal().Err(err).Msg("Error reading payload")
		return "", err
	}

	stringPayload = string(byteData)

	log.Info().Msgf("stringPayload: %s", stringPayload)

	//base64 encode the string
	encodedPayload := base64.StdEncoding.EncodeToString([]byte(stringPayload))

	encodedPayload = routeKey + "." + encodedPayload

	log.Info().Msgf("encodedPayload: %s", encodedPayload)

	return encodedPayload, nil
}

// Returns a stringified version of the response object
func StringifyResponse(response interface{}) (string, error) {

	// stringResponse, err := json.MarshalIndent(response, "", "\t")
	stringResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling response")
		return "", err
	}

	log.Info().Msgf("stringResponse: %s", string(stringResponse))

	return string(stringResponse), nil
}

func SaveToCache(r *http.Request, response interface{}) (string, error) {

	// Prepare the cache key
	routeKey, err := PrepareRouteKey(r)

	if err != nil {
		log.Fatal().Err(err).Msg("Error preparing route key")
		return "", err
	}

	// Prepare the cache key
	cacheKey, err := PrepareCacheKey(r.Body, routeKey)

	if err != nil {
		log.Fatal().Err(err).Msg("Error preparing cache key")
		return "", err
	}

	// Stringify the response
	stringResponse, err := StringifyResponse(response)

	if err != nil {
		log.Fatal().Err(err).Msg("Error stringifying response")
		return "", err
	}

	// Save to cache
	err = redis.SetCache(cacheKey, stringResponse, 0)

	if err != nil {
		log.Fatal().Err(err).Msg("Error saving to cache")
		return "", err
	}

	//TODO: REMOVE LATER
	// CachedResponseToJSON(cacheKey)

	return cacheKey, nil
}

// func CachedResponseToJSON(cacheKey string) (interface{}, error) {
func CachedResponseToJSON(cacheKey string) ([]map[string]interface{}, error) {

	// Get from cache
	cachedResponse, err := redis.GetCache(cacheKey)

	if err != nil {
		log.Fatal().Err(err).Msg("Error getting from cache")
		return nil, err
	}

	log.Info().Msgf("cachedResponse: %s", cachedResponse)

	var cachedResponseJSON []map[string]interface{}
	// var cachedResponseJSON interface{} 

	err = json.Unmarshal([]byte(cachedResponse), &cachedResponseJSON)

	if err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling cached response")
		return nil, err
	}

	log.Info().Msgf("cachedResponseJSON: %v", cachedResponseJSON)
	// log.Info().Msgf("cachedResponseJSON: %t", cachedResponseJSON)

	return cachedResponseJSON, nil
}
