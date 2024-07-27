package pkg

import (
	"github.com/rs/zerolog/log"
)

func CheckAndUpdateIp() error {
	log.Info().Msg("Getting ip")
	ip, err := GetCurrentIpViaCloudflare()

	if err != nil {
		log.Warn().Msg("Could not get ip")
		return err
	}

	log.Info().Str("ip", ip).Msg("Retrieved ip")

	_ = ReconcileState(ip)
	log.Info().Msg("Reconciled; ready to update cloudflare")
	UpdateCloudflare(ip)
	// if result.updated {

	// } else {
	// 	log.Info().Msg("Reconciled; no changes necessary")
	// }

	return nil
}
