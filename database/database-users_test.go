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
