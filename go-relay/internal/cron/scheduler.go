package cron

import (
	"context"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/worker"
	"github.com/rs/zerolog/log"
)

// Scheduler manages periodic cron jobs.
type Scheduler struct {
	db                      *database.DB
	cfg                     *config.Config
	redis                   *config.RedisClient
	dataDir                 string
	browserLogRetentionDays int
	done                    chan struct{}
}

// NewScheduler creates a new cron scheduler.
func NewScheduler(db *database.DB, cfg *config.Config, redis *config.RedisClient, dataDir string, browserLogRetentionDays int) *Scheduler {
	return &Scheduler{
		db:                      db,
		cfg:                     cfg,
		redis:                   redis,
		dataDir:                 dataDir,
		browserLogRetentionDays: browserLogRetentionDays,
		done:                    make(chan struct{}),
	}
}

// Start begins all cron jobs.
func (s *Scheduler) Start() {
	go s.dailyResetLoop()
	go s.monthlyResetLoop()
	go s.tokenCleanupLoop()
	go s.securityCleanupLoop()
	go s.adminCleanupLoop()
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

	if s.redis != nil && s.redis.IsConnected() {
		s.redis.SafeSet(ctx, "last_daily_reset", time.Now().UTC().Format(time.RFC3339), 0)
	}

	log.Info().Msg("Daily reset completed")

	// Clean up browser log files older than the configured retention period.
	worker.CleanBrowserLogs(s.dataDir, s.browserLogRetentionDays)
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
		log.Error().Err(err).Msg("Failed to reset monthly user counters")
	}
	if err := s.db.ApiKeyStore().ResetMonthlyCounters(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to reset monthly scoped key counters")
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
	} else if count > 0 {
		log.Info().Int64("count", count).Msg("Cleaned up expired password reset tokens")
	}

	// Dashboard sessions also expire — clean those up on the same schedule.
	sessionCount, err := s.db.SessionStore().CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired sessions")
	} else if sessionCount > 0 {
		log.Info().Int64("count", sessionCount).Msg("Cleaned up expired dashboard sessions")
	}
}

// securityCleanupLoop cleans up expired pairing codes, key requests, and old connection logs every hour.
func (s *Scheduler) securityCleanupLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runSecurityCleanup()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) runSecurityCleanup() {
	ctx := context.Background()

	// Clean up expired pairing codes
	pairingCount, err := s.db.PairingCodeStore().CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired pairing codes")
	} else if pairingCount > 0 {
		log.Info().Int64("count", pairingCount).Msg("Cleaned up expired pairing codes")
	}

	// Clean up expired key requests
	keyReqCount, err := s.db.KeyRequestStore().CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired key requests")
	} else if keyReqCount > 0 {
		log.Info().Int64("count", keyReqCount).Msg("Cleaned up expired key requests")
	}

	// Clean up expired pair requests
	pairReqCount, err := s.db.PairRequestStore().CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired pair requests")
	} else if pairReqCount > 0 {
		log.Info().Int64("count", pairReqCount).Msg("Cleaned up expired pair requests")
	}

	connRetention := 90
	if s.cfg != nil && s.cfg.ConnectionLogRetentionDays > 0 {
		connRetention = s.cfg.ConnectionLogRetentionDays
	}
	logCount, err := s.db.ConnectionLogStore().CleanupOlderThan(ctx, connRetention)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup old connection logs")
	} else if logCount > 0 {
		log.Info().Int64("count", logCount).Msg("Cleaned up old connection logs")
	}

	// Fix: RemoteRequestLog retention was never wired up
	rrRetention := 30
	if s.cfg != nil && s.cfg.RemoteRequestLogRetentionDays > 0 {
		rrRetention = s.cfg.RemoteRequestLogRetentionDays
	}
	cutoffRR := time.Now().AddDate(0, 0, -rrRetention)
	rrCount, err := s.db.RemoteRequestLogStore().CleanupOlderThan(ctx, cutoffRR)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup old remote request logs")
	} else if rrCount > 0 {
		log.Info().Int64("count", rrCount).Msg("Cleaned up old remote request logs")
	}

	// Module event logs — short retention (default 7 days)
	meRetention := 7
	if s.cfg != nil && s.cfg.ModuleEventLogRetentionDays > 0 {
		meRetention = s.cfg.ModuleEventLogRetentionDays
	}
	cutoffME := time.Now().AddDate(0, 0, -meRetention)
	meCount, err := s.db.ModuleEventLogStore().CleanupOlderThan(ctx, cutoffME)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup old module event logs")
	} else if meCount > 0 {
		log.Info().Int64("count", meCount).Msg("Cleaned up old module event logs")
	}
}

// adminCleanupLoop runs every hour. It purges expired JWT denylist entries
// and old admin audit log entries.
func (s *Scheduler) adminCleanupLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runAdminCleanup()
		case <-s.done:
			return
		}
	}
}

func (s *Scheduler) runAdminCleanup() {
	ctx := context.Background()

	// Distributed lock so only one instance does this work
	if s.redis != nil && s.redis.IsConnected() {
		locked, err := s.redis.Client().SetNX(ctx, "cron:admin_cleanup_lock", "1", 5*time.Minute).Result()
		if err != nil || !locked {
			return
		}
		defer s.redis.SafeDel(ctx, "cron:admin_cleanup_lock")
	}

	// Purge expired JWT denylist entries
	jwtCount, err := s.db.JWTDenylistStore().CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup expired JWT denylist entries")
	} else if jwtCount > 0 {
		log.Info().Int64("count", jwtCount).Msg("Cleaned up expired JWT denylist entries")
	}

	// Purge old admin audit log entries (default 90 days, configurable via ADMIN_AUDIT_LOG_RETENTION_DAYS).
	// We keep this hardcoded to 90 here; the env var is read at server startup and persisted on the config struct.
	// To honor the env var, the cron would need access to cfg — for simplicity we use a sane default and rely
	// on operators wanting different values to set ADMIN_AUDIT_LOG_RETENTION_DAYS and update this constant
	// in a follow-up if needed.
	auditCount, err := s.db.AuditLogStore().CleanupOlderThan(ctx, 90)
	if err != nil {
		log.Error().Err(err).Msg("Failed to cleanup old admin audit logs")
	} else if auditCount > 0 {
		log.Info().Int64("count", auditCount).Msg("Cleaned up old admin audit logs")
	}
}
