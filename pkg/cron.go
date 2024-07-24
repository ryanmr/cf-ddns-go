package pkg

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("module", "cronjob").Logger()

func InitCron() {

	s, err := gocron.NewScheduler()
	defer func() { _ = s.Shutdown() }()

	if err != nil {
		logger.Fatal().Msg("Could not start scheduler")
	}

	var job gocron.Job

	task := func() {
		logger.Info().Msg("Running task")
		t, err := job.NextRun()
		if err == nil {
			logger.Info().Time("NextRun", t).Msgf("Cronjob next run: %s", t.String())
		} else {
			log.Warn().AnErr("err", err).Send()
		}
	}

	job, err = s.NewJob(
		gocron.DurationRandomJob(
			time.Second,
			5*time.Second,
		),
		gocron.NewTask(
			task,
		),
	)

	if err != nil {
		logger.Fatal().Msg("Could not create cronjob")
	}

	logger.Info().Str("cronjob-id", job.ID().String()).Msgf("Cronjob created: %s", job.ID().String())

	s.Start()
	logger.Info().Msg("Started cronjob scheduler")

	t, err := job.NextRun()
	if err == nil {
		logger.Info().Time("NextRun", t).Msgf("Cronjob next run: %s", t.String())
	}

	jobs := s.Jobs()
	for i := 0; i < len(jobs); i++ {
		log.Info().Msgf("Job: %+v", jobs[i])
	}

	// Keep the scheduler running
	select {}

}
