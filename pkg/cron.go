package pkg

import (
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

func InitCron() {

	logger := log.With().Str("module", "cronjob").Logger()

	frequency := 3
	jitterEnv, ok := os.LookupEnv("CRON_FREQUENCY")
	if ok {
		s, err := strconv.Atoi(jitterEnv)
		if err != nil {
			logger.Warn().Msg("Could not evaluate CRON_FREQUENCY, defaulting")
		} else {
			frequency = s
		}
	}

	jitter := 1
	jitterEnv, ok = os.LookupEnv("CRON_FREQUENCY_JITTER")
	if ok {
		s, err := strconv.Atoi(jitterEnv)
		if err != nil {
			logger.Warn().Msg("Could not evaluate CRON_FREQUENCY_JITTER, defaulting")
		} else {
			frequency = s
		}
	}

	s, err := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	if err != nil {
		logger.Fatal().Msg("Could not start scheduler")
	}

	var job gocron.Job

	task := func() {
		logger.Info().Msg("Running task")

		CheckAndUpdateIp()

		logger.Info().Msg("Completed task")
		t, err := job.NextRun()
		if err == nil {
			logger.Info().Time("NextRun", t).Msgf("Cronjob next run: %s", t.String())
		} else {
			log.Warn().AnErr("err", err).Send()
		}
	}

	job, err = s.NewJob(
		gocron.DurationRandomJob(
			(time.Duration(frequency-jitter)*time.Minute),
			time.Duration(frequency+jitter)*time.Minute,
		),
		gocron.NewTask(
			task,
		),
	)

	if err != nil {
		logger.Fatal().Msg("Could not create cronjob")
	}

	logger.Info().
		Str("cronjob-id", job.ID().String()).
		Int("frequency", frequency).
		Int("jitter", jitter).
		Msgf("Cronjob created: %s", job.ID().String())

	s.Start()
	logger.Info().Msg("Started cronjob scheduler")

	t, err := job.NextRun()
	if err == nil {
		logger.Info().Time("NextRun", t).Msgf("Cronjob next run: %s", t.String())
	}

	// keep the scheduler running
	select {}
}
