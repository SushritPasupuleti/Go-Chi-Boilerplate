// Contains the configuration for the logger.
// Utilizes the `zerolog` package.
package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/joho/godotenv"
)

// Initializes the logger with enironment-specific settings
func InitLogging() {
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger() // Enable Caller Info

	err := godotenv.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error loading .env file")
	}

	ENVIRONMENT := os.Getenv("ENVIRONMENT")
	if ENVIRONMENT == "" {
		log.Fatal().
			Err(err).
			Msg("$ENVIRONMENT must be set")
	}

	// Use Pretty Console Logging in Development
	if ENVIRONMENT == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Debug().Msg("Debug logging enabled")

	log.Info().Msg("Info logging enabled")
}
