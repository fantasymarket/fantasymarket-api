package database

import (
	"errors"
	"fantasymarket/database/models"
	"fantasymarket/utils/hash"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	petname "github.com/dustinkirkland/golang-petname"
)

var initialBalance int64 = 1e7

// CreateGuest creates a new Guest User Account
// To lower the barrier of entry to our website,
// this guest account is created automatically
func (s *Service) CreateGuest() (*models.User, error) {

	var username string
	// we generate usenames until the username is unique so everyone
	// starts of with a fresh account
	for username != "" {
		u := petname.Generate(3, "-")

		if s.DB.Where("username = ?", u).Select("username").First(models.User{}).RecordNotFound() {
			username = u
		}
	}

	user := models.User{
		Username: username,
		Portfolio: models.Portfolio{
			Balance: initialBalance,
		},
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &models.User{
		CreatedAt: user.CreatedAt,
		UserID:    user.UserID,
		Username:  user.Username,
	}, nil
}

// ChangePassword changes the password of an existing or
// adds a new password to an previously unprotected user account
// NOTE: this should only be able to be called on your current username
func (s *Service) ChangePassword(username, currentPassword string, newPassword string) error {

	var user models.User
	if err := s.DB.Where(models.User{
		Username: username,
	}).Select("username, password").First(&user).Error; err != nil {
		return err
	}

	if user.Password != "" {
		err := hash.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
		if err != nil {
			return errors.New("current password didn't match")
		}
	}

	passwordHash, err := hash.GeneratePasswordHash([]byte(newPassword))
	if err != nil {
		return err
	}

	return s.DB.Model(&user).Update("password", passwordHash).Error
}

// RenameUser renames a user account
// NOTE: this should only be able to be called on your current username
func (s *Service) RenameUser(userID uuid.UUID, username, newUsername string) error {

	var user models.User
	if err := s.DB.Where(models.User{
		Username: username,
		UserID:   userID,
	}).Select("username. password").First(&user).Error; err != nil {
		return err
	}

	var newUser models.User
	if err := s.DB.Where(models.User{
		Username: newUsername,
	}).Select("username, password").First(&user).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	if newUser.Username != "" {
		return errors.New("username alreay exists")
	}

	if newUser.Password != "" {
		return errors.New("username alreay exists and is protected by a password")
	}

	return s.DB.Model(&user).Update("username", newUsername).Error
}

// LoginUser logs into an account
// NOTE: this should only be called if `RenameUser` fails
func (s *Service) LoginUser(username, password string) (*models.User, error) {

	var user models.User
	if err := s.DB.Where(models.User{
		Username: username,
	}).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("couldn't find user")
		}

		fmt.Println("error loging in:", err)
		return nil, errors.New("could't find user in database")
	}

	if user.Password != "" {
		err := hash.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return nil, errors.New("password didn't match")
		}
	}

	return &models.User{
		CreatedAt: user.CreatedAt,
		UserID:    user.UserID,
		Username:  user.Username,
	}, nil
}

// GetUser searches the db for a user
func (s *Service) GetUser(username string) (*models.User, error) {

	var user models.User
	if err := s.DB.Where(models.User{
		Username: username,
	}).Preload("Portfolio.Items").First(&user).Error; err != nil {
		return nil, errors.New("couldn't find user")
	}

	return &models.User{
		CreatedAt: user.CreatedAt,
		UserID:    user.UserID,
		Username:  user.Username,
		Portfolio: models.Portfolio{
			Balance: user.Portfolio.Balance,
			Items:   user.Portfolio.Items,
		},
	}, nil
}

// GetSelf searches the db for a user (includes private information`)
func (s *Service) GetSelf(username string) (*models.User, error) {

	var user models.User
	if err := s.DB.Where(models.User{
		Username: username,
	}).Preload("Portfolio.Items").First(&user).Error; err != nil {
		return nil, errors.New("couldn't find user")
	}

	return &models.User{
		CreatedAt: user.CreatedAt,
		UserID:    user.UserID,
		Username:  user.Username,
		Portfolio: models.Portfolio{
			Balance: user.Portfolio.Balance,
			Items:   user.Portfolio.Items,
		},
	}, nil
}
