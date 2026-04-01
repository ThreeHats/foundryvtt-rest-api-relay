package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
)

// MemoryStore provides in-memory storage for development mode.
type MemoryStore struct {
	users  *MemoryUserStore
	keys   *MemoryApiKeyStore
	tokens *MemoryPasswordResetTokenStore
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:  &MemoryUserStore{byEmail: make(map[string]*model.User), byAPIKey: make(map[string]*model.User)},
		keys:   &MemoryApiKeyStore{byKey: make(map[string]*model.ApiKey), byUser: make(map[int64][]*model.ApiKey)},
		tokens: &MemoryPasswordResetTokenStore{byID: make(map[int64]*model.PasswordResetToken)},
	}
}

func (m *MemoryStore) Users() model.UserStore                          { return m.users }
func (m *MemoryStore) ApiKeys() model.ApiKeyStore                      { return m.keys }
func (m *MemoryStore) PasswordResetTokens() model.PasswordResetTokenStore { return m.tokens }

// --- Memory User Store ---

type MemoryUserStore struct {
	mu       sync.RWMutex
	byEmail  map[string]*model.User
	byAPIKey map[string]*model.User
	nextID   int64
}

func (s *MemoryUserStore) FindByID(ctx context.Context, id int64) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (s *MemoryUserStore) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u := s.byEmail[email]
	return u, nil
}

func (s *MemoryUserStore) FindByAPIKey(ctx context.Context, apiKey string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u := s.byAPIKey[apiKey]
	return u, nil
}

func (s *MemoryUserStore) FindByStripeCustomerID(ctx context.Context, customerID string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.byEmail {
		if u.StripeCustomerID.Valid && u.StripeCustomerID.String == customerID {
			return u, nil
		}
	}
	return nil, nil
}

func (s *MemoryUserStore) Create(ctx context.Context, user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byEmail[user.Email]; exists {
		return fmt.Errorf("user already exists")
	}
	s.nextID++
	user.ID = s.nextID
	now := time.Now()
	user.CreatedAt = model.NewSQLiteTime(now)
	user.UpdatedAt = model.NewSQLiteTime(now)
	s.byEmail[user.Email] = user
	s.byAPIKey[user.APIKey] = user
	return nil
}

func (s *MemoryUserStore) Update(ctx context.Context, user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.UpdatedAt = model.NewSQLiteTime(time.Now())
	s.byEmail[user.Email] = user
	s.byAPIKey[user.APIKey] = user
	return nil
}

func (s *MemoryUserStore) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for email, u := range s.byEmail {
		if u.ID == id {
			delete(s.byEmail, email)
			delete(s.byAPIKey, u.APIKey)
			return nil
		}
	}
	return nil
}

func (s *MemoryUserStore) IncrementRequests(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.byEmail {
		if u.ID == id {
			u.RequestsThisMonth++
			u.RequestsToday++
			now := time.Now()
			u.LastRequestDate = &model.SQLiteTime{Time: now, Valid: true}
			return nil
		}
	}
	return nil
}

func (s *MemoryUserStore) ResetMonthlyRequests(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.byEmail {
		u.RequestsThisMonth = 0
		u.RequestsToday = 0
	}
	return nil
}

func (s *MemoryUserStore) ResetDailyRequests(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.byEmail {
		u.RequestsToday = 0
	}
	return nil
}

func (s *MemoryUserStore) FindAll(ctx context.Context) ([]*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]*model.User, 0, len(s.byEmail))
	for _, u := range s.byEmail {
		users = append(users, u)
	}
	return users, nil
}

// --- Memory ApiKey Store ---

type MemoryApiKeyStore struct {
	mu     sync.RWMutex
	byKey  map[string]*model.ApiKey
	byUser map[int64][]*model.ApiKey
	nextID int64
}

func (s *MemoryApiKeyStore) FindByKey(ctx context.Context, key string) (*model.ApiKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.byKey[key], nil
}

func (s *MemoryApiKeyStore) FindByID(ctx context.Context, id int64) (*model.ApiKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, k := range s.byKey {
		if k.ID == id {
			return k, nil
		}
	}
	return nil, nil
}

func (s *MemoryApiKeyStore) FindAllByUser(ctx context.Context, userID int64) ([]*model.ApiKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := s.byUser[userID]
	if keys == nil {
		return []*model.ApiKey{}, nil
	}
	return keys, nil
}

func (s *MemoryApiKeyStore) Create(ctx context.Context, key *model.ApiKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	key.ID = s.nextID
	now := time.Now()
	key.CreatedAt = model.NewSQLiteTime(now)
	key.UpdatedAt = model.NewSQLiteTime(now)
	s.byKey[key.Key] = key
	s.byUser[key.UserID] = append(s.byUser[key.UserID], key)
	return nil
}

func (s *MemoryApiKeyStore) Update(ctx context.Context, key *model.ApiKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key.UpdatedAt = model.NewSQLiteTime(time.Now())
	s.byKey[key.Key] = key
	return nil
}

func (s *MemoryApiKeyStore) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.byKey {
		if v.ID == id {
			delete(s.byKey, k)
			userKeys := s.byUser[v.UserID]
			for i, uk := range userKeys {
				if uk.ID == id {
					s.byUser[v.UserID] = append(userKeys[:i], userKeys[i+1:]...)
					break
				}
			}
			return nil
		}
	}
	return nil
}

func (s *MemoryApiKeyStore) DeleteAllByUser(ctx context.Context, userID int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := s.byUser[userID]
	count := int64(len(keys))
	for _, k := range keys {
		delete(s.byKey, k.Key)
	}
	delete(s.byUser, userID)
	return count, nil
}

func (s *MemoryApiKeyStore) ResetDailyCounters(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range s.byKey {
		k.RequestsToday = 0
		k.LastRequestDate = nil
	}
	return nil
}

func (s *MemoryApiKeyStore) IncrementRequests(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range s.byKey {
		if k.ID == id {
			k.RequestsToday++
			k.LastRequestDate = &model.SQLiteTime{Time: time.Now(), Valid: true}
			return nil
		}
	}
	return nil
}

// --- Memory PasswordResetToken Store ---

type MemoryPasswordResetTokenStore struct {
	mu     sync.RWMutex
	byID   map[int64]*model.PasswordResetToken
	nextID int64
}

func (s *MemoryPasswordResetTokenStore) FindByTokenHash(ctx context.Context, hash string) (*model.PasswordResetToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.byID {
		if t.TokenHash == hash && !t.Used && t.ExpiresAt.After(time.Now()) {
			return t, nil
		}
	}
	return nil, nil
}

func (s *MemoryPasswordResetTokenStore) Create(ctx context.Context, token *model.PasswordResetToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	token.ID = s.nextID
	now := time.Now()
	token.CreatedAt = now
	token.UpdatedAt = now
	s.byID[token.ID] = token
	return nil
}

func (s *MemoryPasswordResetTokenStore) MarkUsed(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if t, ok := s.byID[id]; ok {
		t.Used = true
		t.UpdatedAt = time.Now()
	}
	return nil
}

func (s *MemoryPasswordResetTokenStore) InvalidateForUser(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.byID {
		if t.UserID == userID && !t.Used {
			t.Used = true
			t.UpdatedAt = time.Now()
		}
	}
	return nil
}

func (s *MemoryPasswordResetTokenStore) CleanupExpired(ctx context.Context) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := time.Now().Add(-24 * time.Hour)
	var count int64
	for id, t := range s.byID {
		if (t.Used || t.ExpiresAt.Before(time.Now())) && t.CreatedAt.Before(cutoff) {
			delete(s.byID, id)
			count++
		}
	}
	return count, nil
}
