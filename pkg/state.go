package pkg

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type ServerState struct {
	PreviousIp         *string
	PreviousUpdateTime *time.Time
	Mutex              sync.Mutex
}

var state ServerState

func InitState() {
	ip := "0.0.0.0"
	time := time.Now()

	log.Debug().Str("ip", ip).Time("time", time).Int8("phase", 1).Msg("Initializing state")

	state.Mutex.Lock()
	state.PreviousIp = &ip
	state.PreviousUpdateTime = &time
	state.Mutex.Unlock()

	log.Debug().Str("ip", ip).Time("time", time).Int8("phase", 2).Msg("Initialized state")

}

func UpdateState(ip string, time time.Time) {

	log.Debug().Str("ip", ip).Time("time", time).Int8("phase", 1).Msg("Updating state")

	state.Mutex.Lock()
	*state.PreviousIp = ip
	*state.PreviousUpdateTime = time
	state.Mutex.Unlock()

	log.Debug().Str("ip", ip).Time("time", time).Int8("phase", 2).Msg("Updated state")

}
