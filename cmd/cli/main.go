package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ryanmr/cf-ddns-go/pkg"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set log level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log format
	switch strings.ToLower(os.Getenv("LOG_FORMAT")) {
	case "plain":

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	case "color":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
	case "json":
		// Default format is JSON, no need to change
	default:
		// Default format is JSON, no need to change
	}

	log.Info().Msg("cf-ddns-go")
	log.Info().Msg("repo: https://github.com/ryanmr/cf-ddns-go?from=src")

	pkg.Serve()
}
