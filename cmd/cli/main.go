package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ryanmr/cf-ddns-go/pkg"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("cf-ddns-go")
	log.Info().Msg("repo: https://github.com/ryanmr/cf-ddns-go?from=src")

	pkg.Serve()
}
