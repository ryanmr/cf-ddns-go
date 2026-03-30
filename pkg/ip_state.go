package pkg

import (
	"github.com/rs/zerolog/log"
)

func CheckAndUpdateIp(override bool) error {
	log.Info().Msg("Getting ip")
	ip, err := GetCurrentIpViaCloudflare()

	if err != nil {
		log.Warn().Msg("Could not get ip")
		return err
	}

	log.Info().Str("ip", ip).Msg("Retrieved ip")

	result := ReconcileState(ip)

	if result.updated {
		handleDiscordWebhook(ip)
	}

	if override {
		log.Info().Msg("Reconciliation override; updating cloudflare anyway")
		UpdateCloudflare(ip)
	} else if result.updated {
		log.Info().Msg("Reconciled; ip changed, updating cloudflare")
		UpdateCloudflare(ip)
	} else {
		// Even if in-memory state hasn't changed, verify Cloudflare has the right IP.
		// This catches stale DNS records from manual edits, container restarts, etc.
		dnsIp, err := GetCloudflareDnsIp()
		if err != nil {
			log.Warn().Err(err).Msg("Could not verify Cloudflare DNS record")
		} else if dnsIp != ip {
			log.Info().
				Str("dns-ip", dnsIp).
				Str("current-ip", ip).
				Msg("Cloudflare DNS record is stale; updating")
			UpdateCloudflare(ip)
			handleDiscordWebhook(ip)
		} else {
			log.Info().Msg("Reconciled; no changes necessary")
		}
	}

	return nil
}
