package repositories

import (
	"auth-app/models"
	"database/sql"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByID(id string) (*models.User, error)
	Delete(id string) (int64, error)
	Update(id string, user *models.User) (int64, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	err := r.db.QueryRow(
		"INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id",
		user.Email,
		user.Name,
		user.Password,
	).Scan(&user.ID)

	return err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, email FROM users WHERE id=$1", id).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *userRepository) Update(id string, user *models.User) (int64, error) {
	result, err := r.db.Exec("UPDATE users SET email=$1, name=$2 WHERE id=$3", user.Email, user.Name, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
