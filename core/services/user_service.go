package services

import (
	"database/sql"
	"dsn/core/types"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(req types.CreateUserRequest) (*types.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// if no users exist yet, make this user an admin
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, err
	}
	isAdmin := count == 0

	query := `INSERT INTO users (username, email, password_hash, is_admin) 
		VALUES (?, ?, ?, ?)
		RETURNING id, created_at, updated_at`

	var user types.User
	err = s.db.QueryRow(query, req.Username, req.Email, string(hashedPassword), isAdmin).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.IsAdmin = isAdmin

	return &user, nil
}

func (s *UserService) GetByUsername(username string) (*types.User, error) {
	query := `SELECT id, username, email, password_hash, is_admin, created_at, updated_at 
		FROM users 
		WHERE username = ?`

	var user types.User
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetByID(id int) (*types.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`

	var user types.User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves all users (admin only)
func (s *UserService) GetAll() ([]types.User, error) {
	query := `
		SELECT id, username, email, is_admin, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email,
			&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func (s *UserService) ValidatePassword(user *types.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
