package pkg

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type ServerState struct {
	CurrentIp  *string
	PreviousIp *string
	UpdateTime *time.Time
	Mutex      sync.Mutex
}

var state ServerState

// Initializes the server state.
//
// This must be called before any calls to UpdateState.
func InitState() {
	ip, err := GetCurrentIpViaCloudflare()
	if err != nil {
		log.Warn().Msg("Could not get ip during InitState")
		return
	}

	currentIp := ip + ""
	previousIp := ""
	time := time.Now()

	log.Debug().Msg("Initializing state")

	state.Mutex.Lock()
	state.CurrentIp = &currentIp
	state.PreviousIp = &previousIp
	state.UpdateTime = &time
	state.Mutex.Unlock()

	log.Debug().Msg("Initialized state")

}

type ReconcileResult struct {
	updated bool
}

// Reconcile the server state.
func ReconcileState(ip string) ReconcileResult {
	time := time.Now()

	log.Info().
		Msg("Updating state")

	state.Mutex.Lock()

	log.Debug().
		Str("new-ip", ip).
		Str("current-ip", *state.CurrentIp).
		Str("previous-ip", *state.PreviousIp).
		Str("update-ip", (*state.UpdateTime).String()).
		Msg("Reviewing state conditions")

	updated := false
	if ip != *state.CurrentIp {
		*state.PreviousIp = *state.CurrentIp
		*state.UpdateTime = time
		*state.CurrentIp = ip

		updated = true

		log.Debug().
			Str("new-ip", ip).
			Str("current-ip", *state.CurrentIp).
			Str("previous-ip", *state.PreviousIp).
			Str("update-ip", (*state.UpdateTime).String()).
			Msg("Realizing state condition")
	}

	state.Mutex.Unlock()

	log.Info().
		Msg("Updated state")

	result := ReconcileResult{updated}
	return result
}
