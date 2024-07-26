package pkg

import (
	"github.com/rs/zerolog/log"
)

func UpdateCloudflare(ip string) error {
	log.Info().Msg("Updating cloudflare dns settings")
	return nil
}
