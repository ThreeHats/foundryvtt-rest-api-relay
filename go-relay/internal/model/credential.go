package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Credential represents stored Foundry VTT credentials for a user.
type Credential struct {
	ID                       int64      `db:"id" json:"id"`
	UserID                   int64      `db:"userId" json:"userId"`
	Name                     string     `db:"name" json:"name"`
	FoundryURL               string     `db:"foundryUrl" json:"foundryUrl"`
	FoundryUsername           string     `db:"foundryUsername" json:"foundryUsername"`
	EncryptedFoundryPassword string     `db:"encryptedFoundryPassword" json:"-"`
	PasswordIV               string     `db:"passwordIv" json:"-"`
	PasswordAuthTag          string     `db:"passwordAuthTag" json:"-"`
	CreatedAt                SQLiteTime `db:"createdAt" json:"createdAt"`
	UpdatedAt                SQLiteTime `db:"updatedAt" json:"updatedAt"`
}

// CredentialStore defines operations on stored Foundry credentials.
type CredentialStore interface {
	FindByID(ctx context.Context, id int64) (*Credential, error)
	FindAllByUser(ctx context.Context, userID int64) ([]*Credential, error)
	Create(ctx context.Context, cred *Credential) error
	Update(ctx context.Context, cred *Credential) error
	Delete(ctx context.Context, id int64) error
}

// SQLCredentialStore implements CredentialStore with sqlx.
type SQLCredentialStore struct {
	DB     DBTX
	DBType string
}

func (s *SQLCredentialStore) tableName() string {
	if s.DBType == "sqlite" {
		return "Credentials"
	}
	return `"Credentials"`
}

func (s *SQLCredentialStore) col(name string) string {
	return Col(s.DBType, name)
}

// errCorruptCredential is returned when a credential record is missing
// one or more required encryption fields (IV or auth tag). This indicates
// database corruption or a partial write.
var errCorruptCredential = errors.New("credential record is missing required encryption fields")

func (s *SQLCredentialStore) FindByID(ctx context.Context, id int64) (*Credential, error) {
	var c Credential
	err := s.DB.GetContext(ctx, &c, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName()), id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if c.EncryptedFoundryPassword != "" && (c.PasswordIV == "" || c.PasswordAuthTag == "") {
		return nil, errCorruptCredential
	}
	return &c, nil
}

func (s *SQLCredentialStore) FindAllByUser(ctx context.Context, userID int64) ([]*Credential, error) {
	var creds []*Credential
	err := s.DB.SelectContext(ctx, &creds, fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", s.tableName(), s.col("user_id")), userID)
	if creds == nil {
		creds = []*Credential{}
	}
	return creds, err
}

func (s *SQLCredentialStore) Create(ctx context.Context, cred *Credential) error {
	now := time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (%s, name, %s, %s, %s, %s, %s, %s, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		s.tableName(),
		s.col("user_id"), s.col("foundry_url"), s.col("foundry_username"),
		s.col("encrypted_foundry_password"), s.col("password_iv"), s.col("password_auth_tag"),
		s.col("created_at"), s.col("updated_at"))

	if s.DBType != "sqlite" {
		query += " RETURNING id"
		return s.DB.QueryRowContext(ctx, query,
			cred.UserID, cred.Name, cred.FoundryURL, cred.FoundryUsername,
			cred.EncryptedFoundryPassword, cred.PasswordIV, cred.PasswordAuthTag,
			now, now,
		).Scan(&cred.ID)
	}

	result, err := s.DB.ExecContext(ctx, query,
		cred.UserID, cred.Name, cred.FoundryURL, cred.FoundryUsername,
		cred.EncryptedFoundryPassword, cred.PasswordIV, cred.PasswordAuthTag,
		now, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	cred.ID = id
	cred.CreatedAt = NewSQLiteTime(now)
	cred.UpdatedAt = NewSQLiteTime(now)
	return nil
}

func (s *SQLCredentialStore) Update(ctx context.Context, cred *Credential) error {
	cred.UpdatedAt = NewSQLiteTime(time.Now())
	query := fmt.Sprintf(`UPDATE %s SET name=$1, %s=$2, %s=$3,
		%s=$4, %s=$5, %s=$6, %s=$7
		WHERE id=$8`,
		s.tableName(),
		s.col("foundry_url"), s.col("foundry_username"),
		s.col("encrypted_foundry_password"), s.col("password_iv"), s.col("password_auth_tag"),
		s.col("updated_at"))
	_, err := s.DB.ExecContext(ctx, query,
		cred.Name, cred.FoundryURL, cred.FoundryUsername,
		cred.EncryptedFoundryPassword, cred.PasswordIV, cred.PasswordAuthTag,
		cred.UpdatedAt, cred.ID)
	return err
}

func (s *SQLCredentialStore) Delete(ctx context.Context, id int64) error {
	_, err := s.DB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName()), id)
	return err
}
