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
	currentIp := "0.0.0.0"
	previousIp := "0.0.0.0"
	time := time.Now()

	log.Debug().Msg("Initializing state")

	state.Mutex.Lock()
	state.CurrentIp = &currentIp
	state.PreviousIp = &previousIp
	state.UpdateTime = &time
	state.Mutex.Unlock()

	log.Debug().Msg("Initialized state")

}

// Updates the server state.
func UpdateState(ip string) {
	time := time.Now()

	log.Debug().
		Str("ip", ip).
		Time("time", time).
		Int8("phase", 1).
		Msg("Updating state")

	state.Mutex.Lock()
	*state.PreviousIp = *state.CurrentIp
	*state.UpdateTime = time
	*state.CurrentIp = ip
	state.Mutex.Unlock()

	log.Debug().
		Str("ip", ip).
		Time("time", time).
		Int8("phase", 2).
		Msg("Updated state")

}
