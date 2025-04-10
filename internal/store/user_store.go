package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassowrd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassowrd), 10)
	if err != nil {
		return err
	}
	p.plainText = &plainTextPassowrd
	p.hash = hash
	return nil
}

func (p *password) isMatch(plainTextPassowrd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassowrd))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUserName(id int) (*User, error)
	UpdateUser(*User) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (user_name,email,password_hash,bio)
	VALUES ($1,$2,$3,$4)
	RETURNING id,created_at,updated_at
	`
	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (pg *PostgresUserStore) GetUserByUserName(id int) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	SELECT id,username,email,password_hash,bio,created_at,updated_at
	FROM users 
	WHERE username=$1
	`
	err := pg.db.QueryRow(query, user.Username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (pg *PostgresUserStore) UpdateUser(user *User) error {

	query := `
	UPDATE useres
	SET username=$1,email=$2,bio=$3,updated_at=CURRENT_TIMESTAMP
	WHERE id =$4 
	`
	res, err := pg.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
