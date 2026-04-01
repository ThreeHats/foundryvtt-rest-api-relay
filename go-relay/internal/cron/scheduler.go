package cron

import (
	"context"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/rs/zerolog/log"
)

// Scheduler manages periodic cron jobs.
type Scheduler struct {
	db    *database.DB
	redis *config.RedisClient
	done  chan struct{}
}

// NewScheduler creates a new cron scheduler.
func NewScheduler(db *database.DB, redis *config.RedisClient) *Scheduler {
	return &Scheduler{
		db:    db,
		redis: redis,
		done:  make(chan struct{}),
	}
}

// Start begins all cron jobs.
func (s *Scheduler) Start() {
	go s.dailyResetLoop()
	go s.monthlyResetLoop()
	go s.tokenCleanupLoop()
	log.Info().Msg("Cron jobs started")
}

// Stop gracefully stops all cron jobs.
func (s *Scheduler) Stop() {
	close(s.done)
}

// dailyResetLoop resets daily request counters at midnight UTC.
func (s *Scheduler) dailyResetLoop() {
	for {
		now := time.Now().UTC()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		timer := time.NewTimer(nextMidnight.Sub(now))

		select {
		case <-timer.C:
			s.runDailyReset()
		case <-s.done:
			timer.Stop()
			return
		}
	}
}

// monthlyResetLoop resets monthly request counters on the 1st of each month.
func (s *Scheduler) monthlyResetLoop() {
	for {
		now := time.Now().UTC()
		nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		timer := time.NewTimer(nextMonth.Sub(now))

		select {
		case <-timer.C:
			s.runMonthlyReset()
		case <-s.done:
			timer.Stop()
			return
		}
	}
}

// tokenCleanupLoop cleans up expired password reset tokens every 6 hours.
func (s *Scheduler) tokenCleanupLoop() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runTokenCleanup()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) runDailyReset() {
	ctx := context.Background()

	// Distributed lock via Redis
	if s.redis != nil && s.redis.IsConnected() {
		locked, err := s.redis.Client().SetNX(ctx, "cron:daily_reset_lock", "1", 5*time.Minute).Result()
		if err != nil || !locked {
			log.Debug().Msg("Daily reset already running on another instance")
			return
		}
		defer s.redis.SafeDel(ctx, "cron:daily_reset_lock")
	}

	log.Info().Msg("Running daily request counter reset")

	if err := s.db.UserStore().ResetDailyRequests(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to reset daily user counters")
	}
	if err := s.db.ApiKeyStore().ResetDailyCounters(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to reset daily API key counters")
	}

	if s.redis != nil && s.redis.IsConnected() {
		s.redis.SafeSet(ctx, "last_daily_reset", time.Now().UTC().Format(time.RFC3339), 0)
	}

	log.Info().Msg("Daily reset completed")
}

func (s *Scheduler) runMonthlyReset() {
	ctx := context.Background()

	if s.redis != nil && s.redis.IsConnected() {
		locked, err := s.redis.Client().SetNX(ctx, "cron:monthly_reset_lock", "1", 5*time.Minute).Result()
		if err != nil || !locked {
			log.Debug().Msg("Monthly reset already running on another instance")
			return
		}
		defer s.redis.SafeDel(ctx, "cron:monthly_reset_lock")
	}

	log.Info().Msg("Running monthly request counter reset")

	if err := s.db.UserStore().ResetMonthlyRequests(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to reset monthly counters")
	}

	if s.redis != nil && s.redis.IsConnected() {
		s.redis.SafeSet(ctx, "last_monthly_reset", time.Now().UTC().Format(time.RFC3339), 0)
	}

	log.Info().Msg("Monthly reset completed")
}

func (s *Scheduler) runTokenCleanup() {
	ctx := context.Background()
	count, err := s.db.PasswordResetTokenStore().CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired tokens")
		return
	}
	if count > 0 {
		log.Info().Int64("count", count).Msg("Cleaned up expired password reset tokens")
	}
}
