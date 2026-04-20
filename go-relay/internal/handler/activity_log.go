package handler

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/middleware"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
)

// ActivityRouter mounts GET /auth/activity for the authenticated user.
func ActivityRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(db, nil))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil || reqCtx.User == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		// Scoped keys cannot access logs — require master key or session.
		if reqCtx.ScopedKey != nil {
			helpers.WriteError(w, http.StatusForbidden, "Scoped keys cannot access the activity log")
			return
		}

		user, ok := reqCtx.User.(*model.User)
		if !ok || user == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		filters, limit, offset := parseActivityParams(r, user.ID)

		events, total, err := mergeActivityLogs(r.Context(), db, filters, limit, offset)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to fetch activity log")
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"events": events,
			"total":  total,
			"offset": offset,
			"limit":  limit,
		})
	})

	return r
}

func parseActivityParams(r *http.Request, userID int64) (model.ActivityFilters, int, int) {
	q := r.URL.Query()

	limit := 50
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if limit > 200 {
		limit = 200
	}

	offset := 0
	if o := q.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	filters := model.ActivityFilters{
		UserID: userID,
		Type:   q.Get("type"),
		World:  q.Get("world"),
		Action: q.Get("action"),
		Limit:  limit,
		Offset: offset,
	}

	if s := q.Get("success"); s == "true" {
		t := true
		filters.Success = &t
	} else if s == "false" {
		f := false
		filters.Success = &f
	}

	if since := q.Get("since"); since != "" {
		if t, err := time.Parse(time.RFC3339, since); err == nil {
			filters.Since = t
		}
	}
	if until := q.Get("until"); until != "" {
		if t, err := time.Parse(time.RFC3339, until); err == nil {
			filters.Until = t
		}
	}

	return filters, limit, offset
}

// mergeActivityLogs fetches from all three sources (unless filtered), merges,
// sorts by time descending, then slices [offset:offset+limit].
func mergeActivityLogs(ctx context.Context, db *database.DB, filters model.ActivityFilters, limit, offset int) ([]model.ActivityEvent, int64, error) {
	var all []model.ActivityEvent
	var totalConn, totalRR, totalME int64

	fetchFilters := filters
	fetchFilters.Offset = 0 // fetch top (offset+limit) from each source; slice after merge

	if filters.Type == "" || filters.Type == "connection" {
		connLogs, total, err := db.ConnectionLogStore().FindFiltered(ctx, fetchFilters)
		if err != nil {
			return nil, 0, err
		}
		totalConn = total
		for _, c := range connLogs {
			ev := model.ActivityEvent{
				ID:           c.ID,
				Type:         "connection",
				EventSubtype: "connect",
				ClientID:     c.ClientID,
				Flagged:      c.Flagged,
				UserID:       c.UserID,
				CreatedAt:    c.CreatedAt.Time,
			}
			if c.WorldTitle.Valid {
				ev.WorldTitle = c.WorldTitle.String
			}
			if c.TokenName.Valid {
				ev.TokenName = c.TokenName.String
			}
			if c.IPAddress.Valid {
				ev.IPAddress = c.IPAddress.String
			}
			all = append(all, ev)
		}
	}

	if filters.Type == "" || filters.Type == "remote_request" {
		rrLogs, total, err := db.RemoteRequestLogStore().FindFiltered(ctx, fetchFilters)
		if err != nil {
			return nil, 0, err
		}
		totalRR = total
		for _, rr := range rrLogs {
			success := bool(rr.Success)
			ev := model.ActivityEvent{
				ID:             rr.ID,
				Type:           "remote_request",
				ClientID:       rr.SourceClientID,
				TargetClientID: rr.TargetClientID,
				Action:         rr.Action,
				Success:        &success,
				UserID:         rr.UserID,
				CreatedAt:      rr.CreatedAt.Time,
			}
			if rr.ErrorMessage.Valid {
				ev.ErrorMessage = rr.ErrorMessage.String
			}
			all = append(all, ev)
		}
	}

	if filters.Type == "" || filters.Type == "module_event" {
		meLogs, total, err := db.ModuleEventLogStore().FindFiltered(ctx, fetchFilters)
		if err != nil {
			return nil, 0, err
		}
		totalME = total
		for _, me := range meLogs {
			all = append(all, model.ActivityEvent{
				ID:          me.ID,
				Type:        "module_event",
				ClientID:    me.ClientID,
				WorldTitle:  me.WorldTitle,
				Action:      me.EventType,
				Actor:       me.Actor,
				Description: me.Description,
				UserID:      me.UserID,
				CreatedAt:   me.CreatedAt.Time,
			})
		}
	}

	// Sort all events by CreatedAt descending
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	// Apply offset/limit slice
	total := totalConn + totalRR + totalME
	if offset >= len(all) {
		return []model.ActivityEvent{}, total, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], total, nil
}
