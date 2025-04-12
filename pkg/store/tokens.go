package store

import (
	"database/sql"
	"time"

	"github.com/yehiamoh/go-fem-workshop/pkg/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(*tokens.Token) error
	CreateNewToken(userID int, timeToLive time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (pg *PostgresTokenStore) CreateNewToken(userID int, timeToLive time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerteToken(userID, timeToLive, scope)
	if err != nil {
		return nil, err
	}
	err = pg.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (pg *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash,user_id,expiry,scope)
	VALUES($1,$2,$3,$4)
	`
	_, err := pg.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	if err != nil {
		return err
	}
	return nil
}
func (pg *PostgresTokenStore) DeleteAllTokensForUser(userId int, scope string) error {
	query := `
	DELETE FROM tokens 
	WHERE user_id=$1 and scope=$2
	`
	_, err := pg.db.Exec(query, userId, scope)
	if err != nil {
		return err
	}
	return nil
}
