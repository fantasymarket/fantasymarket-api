package database_test

import (
	"fantasymarket/database/models"
	"github.com/stretchr/testify/assert"
)

func (suite *DatabaseTestSuite) TestCreateGuest() {
	for i:=0;i<4;i++ {
		var testUser models.User

		user, err := suite.dbService.CreateGuest()
		assert.Equal(suite.T(), nil, err)

		assert.NotEqual(suite.T(), "", user.Username)

		err = suite.dbService.DB.Where("username = ?", user.Username).First(&testUser).Error
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), user.Username, testUser.Username)
	}
}

// func (s *Service) CreateGuest() (*models.User, error) {

// 	var username string
// 	// we generate usenames until the username is unique so everyone
// 	// starts of with a fresh account
// 	for username == "" {
// 		u := petname.Generate(3, "-")

// 		if s.DB.Where("username = ?", u).Select("username").First(&models.User{}).RecordNotFound() {
// 			username = u
// 		}
// 	}

// 	user := models.User{
// 		Username: username,
// 		Portfolio: models.Portfolio{
// 			Balance: initialBalance,
// 		},
// 	}

// 	if err := s.DB.Create(&user).Error; err != nil {
// 		return nil, err
// 	}

// 	return &models.User{
// 		CreatedAt: user.CreatedAt,
// 		UserID:    user.UserID,
// 		Username:  user.Username,
// 	}, nil
// }