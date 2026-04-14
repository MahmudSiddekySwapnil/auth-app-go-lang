package services

import (
	"auth-app/models"
	"auth-app/repositories"
	"auth-app/utils"
	"database/sql"
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (string, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user *models.User) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User) error {
	// Check if already exists
	_, err := s.repo.FindByEmail(user.Email)
	if err == nil {
		log.Printf("⚠️  Duplicate email error: %s\n", user.Email)
		return errors.New("Email already registered")
	}
	if err != nil && err != sql.ErrNoRows {
		log.Printf("❌ Database error checking email: %v\n", err)
		return errors.New("Database error")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return errors.New("Failed to hash password")
	}
	
	log.Printf("🔐 Password hashed successfully: %s\n", string(hashedPassword))
	user.Password = string(hashedPassword)

	err = s.repo.Create(user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("⚠️  Duplicate email error: %s\n", user.Email)
			return errors.New("Email already registered")
		}
		log.Printf("❌ Database error: %v\n", err)
		return err
	}

	log.Printf("✅ User registered successfully - ID: %d, Email: %s\n", user.ID, user.Email)
	return nil
}

func (s *userService) Login(email, password string) (string, error) {
	log.Printf("🔑 Login attempt - Email: %s\n", email)
	
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("⚠️  User not found: %s\n", email)
			return "", errors.New("Invalid email")
		}
		log.Printf("❌ Database error during login: %v\n", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("⚠️  Wrong password for: %s\n", email)
		return "", errors.New("Wrong password")
	}

	token, err := utils.GenerateToken(email)
	if err != nil {
		log.Printf("❌ Token generation error: %v\n", err)
		return "", errors.New("Failed to generate token")
	}

	log.Printf("✅ Login successful - Email: %s\n", email)
	return token, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	log.Printf("📋 Fetching all users...\n")
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("⚠️  User not found with ID: %s\n", id)
			return nil, errors.New("User not found")
		}
		log.Printf("❌ Database error fetching user by ID: %v\n", err)
		return nil, err
	}
	log.Printf("✅ Fetched user by ID: %d, Email: %s\n", user.ID, user.Email)
	return user, nil
}

func (s *userService) DeleteUser(id string) error {
	rowsAffected, err := s.repo.Delete(id)
	if err != nil {
		log.Printf("❌ Database error deleting user: %v\n", err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("⚠️  No user found to delete with ID: %s\n", id)
		return errors.New("User not found")
	}
	log.Printf("✅ User deleted successfully - ID: %s\n", id)
	return nil
}

func (s *userService) UpdateUser(id string, user *models.User) error {
	rowsAffected, err := s.repo.Update(id, user)
	if err != nil {
		log.Printf("❌ Database error updating user: %v\n", err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("⚠️  No user found to update with ID: %s\n", id)
		return errors.New("User not found")
	}
	log.Printf("✅ User updated successfully - ID: %s, New Email: %s\n", id, user.Email)
	return nil
}
