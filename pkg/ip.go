package pkg

import (
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetCurrentIpViaCloudflare() (string, error) {
	cmd := exec.Command("dig", "@1.1.1.1", "ch", "txt", "whoami.Cloudflare", "+short")
	output, err := cmd.Output()
	if err != nil {
		log.Warn().Msg("Could not fetch ip address")
		return "", err
	}

	// Convert output to string and remove quotes and newlines
	ip := strings.ReplaceAll(string(output), "\"", "")
	ip = strings.TrimSpace(ip)

	return ip, nil
}
