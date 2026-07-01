package model

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// bptr returns a pointer to b — for tri-state filter fields.
func bptr(b bool) *bool { return &b }

func openTestDB(t *testing.T, dbType string) *sqlx.DB {
	t.Helper()
	var db *sqlx.DB
	var err error
	if dbType == "sqlite" {
		path := filepath.Join(t.TempDir(), "filter.db")
		db, err = sqlx.Connect("sqlite", path)
	} else {
		db, err = sqlx.Connect("pgx", os.Getenv("TEST_DATABASE_URL"))
	}
	if err != nil {
		t.Fatalf("connect %s: %v", dbType, err)
	}
	db = db.Unsafe()
	db.MapperFunc(func(s string) string { return strings.ToLower(strings.ReplaceAll(s, "_", "")) })
	t.Cleanup(func() { db.Close() })
	return db
}

func seedUsers(t *testing.T, db *sqlx.DB, dbType string) {
	t.Helper()
	boolType := "INTEGER"
	if dbType == "postgres" {
		boolType = "BOOLEAN"
		db.MustExec(`DROP TABLE IF EXISTS "Users"`)
	}
	create := fmt.Sprintf(`CREATE TABLE "Users" (
		id INTEGER PRIMARY KEY,
		email TEXT,
		role TEXT,
		disabled %[1]s,
		"emailVerified" %[1]s,
		"apiKeyRotationRequired" %[1]s,
		"subscriptionStatus" TEXT,
		"requestsToday" INTEGER,
		"requestsThisMonth" INTEGER,
		"createdAt" TEXT
	)`, boolType)
	db.MustExec(create)

	ins := `INSERT INTO "Users" (id, email, role, disabled, "emailVerified", "apiKeyRotationRequired", "subscriptionStatus", "requestsToday", "requestsThisMonth", "createdAt")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	rows := []struct {
		id          int
		email, role string
		disabled    bool
		verified    bool
		rotation    bool
		sub         interface{}
		today       int
		month       int
		created     string
	}{
		{1, "alice@example.com", "user", false, true, false, "active", 10, 100, "2026-01-01 00:00:00"},
		{2, "bob@example.com", "admin", true, false, true, nil, 5, 50, "2026-01-02 00:00:00"},
		{3, "carol@test.com", "user", false, false, false, "free", 20, 200, "2026-01-03 00:00:00"},
		{4, "dave@example.com", "user", false, true, true, "past_due", 1, 2, "2026-01-04 00:00:00"},
		{5, "admin@example.com", "admin", false, true, false, "active", 0, 0, "2026-01-05 00:00:00"},
	}
	for _, r := range rows {
		db.MustExec(ins, r.id, r.email, r.role, r.disabled, r.verified, r.rotation, r.sub, r.today, r.month, r.created)
	}
}

func gotIDs(users []*User) []int64 {
	out := make([]int64, len(users))
	for i, u := range users {
		out[i] = u.ID
	}
	return out
}

func eqIDs(a []int64, b ...int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func runFilterSuite(t *testing.T, dbType string) {
	db := openTestDB(t, dbType)
	seedUsers(t, db, dbType)
	store := &SQLUserStore{DB: db, DBType: dbType}
	ctx := context.Background()

	cases := []struct {
		name      string
		q         UserQuery
		wantTotal int64
		wantIDs   []int64
	}{
		{"no filter", UserQuery{SortBy: "id"}, 5, []int64{1, 2, 3, 4, 5}},
		{"search email substring", UserQuery{Search: "example", SortBy: "id"}, 4, []int64{1, 2, 4, 5}},
		{"search by id digit", UserQuery{Search: "3", SortBy: "id"}, 1, []int64{3}},
		{"search case-insensitive", UserQuery{Search: "CAROL", SortBy: "id"}, 1, []int64{3}},
		{"role admin", UserQuery{Role: "admin", SortBy: "id"}, 2, []int64{2, 5}},
		{"disabled true", UserQuery{Disabled: bptr(true), SortBy: "id"}, 1, []int64{2}},
		{"disabled false", UserQuery{Disabled: bptr(false), SortBy: "id"}, 4, []int64{1, 3, 4, 5}},
		{"verified true", UserQuery{EmailVerified: bptr(true), SortBy: "id"}, 3, []int64{1, 4, 5}},
		{"verified false", UserQuery{EmailVerified: bptr(false), SortBy: "id"}, 2, []int64{2, 3}},
		{"rotation true", UserQuery{RotationRequired: bptr(true), SortBy: "id"}, 2, []int64{2, 4}},
		{"sub active", UserQuery{Subscription: "active", SortBy: "id"}, 2, []int64{1, 5}},
		{"sub free incl null", UserQuery{Subscription: "free", SortBy: "id"}, 2, []int64{2, 3}},
		{"sub past_due", UserQuery{Subscription: "past_due", SortBy: "id"}, 1, []int64{4}},
		{"sort today asc", UserQuery{SortBy: "requestsToday"}, 5, []int64{5, 4, 2, 1, 3}},
		{"sort today desc", UserQuery{SortBy: "requestsToday", SortDesc: true}, 5, []int64{3, 1, 2, 4, 5}},
		{"sort email asc", UserQuery{SortBy: "email"}, 5, []int64{5, 1, 2, 3, 4}},
		{"sort createdAt desc", UserQuery{SortBy: "createdAt", SortDesc: true}, 5, []int64{5, 4, 3, 2, 1}},
		{"sort subscription", UserQuery{SortBy: "subscriptionStatus"}, 5, nil}, // order varies on NULL; only check total
		{"combined role+verified", UserQuery{Role: "user", EmailVerified: bptr(true), SortBy: "id"}, 2, []int64{1, 4}},
		{"unknown sort falls back to id", UserQuery{SortBy: "evil; DROP TABLE"}, 5, []int64{1, 2, 3, 4, 5}},
		{"pagination page 1", UserQuery{SortBy: "id", Limit: 2, Offset: 0}, 5, []int64{1, 2}},
		{"pagination page 2", UserQuery{SortBy: "id", Limit: 2, Offset: 2}, 5, []int64{3, 4}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			users, total, err := store.FindFilteredPaginated(ctx, tc.q)
			if err != nil {
				t.Fatalf("query error: %v", err)
			}
			if total != tc.wantTotal {
				t.Errorf("total = %d, want %d", total, tc.wantTotal)
			}
			if tc.wantIDs != nil && !eqIDs(gotIDs(users), tc.wantIDs...) {
				t.Errorf("ids = %v, want %v", gotIDs(users), tc.wantIDs)
			}
		})
	}
}

func TestFindFilteredPaginated_SQLite(t *testing.T) {
	runFilterSuite(t, "sqlite")
}

func TestFindFilteredPaginated_Postgres(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping Postgres dialect check")
	}
	runFilterSuite(t, "postgres")
}
