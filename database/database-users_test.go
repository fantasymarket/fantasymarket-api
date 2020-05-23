package database_test

import (
	"fantasymarket/database/models"
	"fantasymarket/utils/hash"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *DatabaseTestSuite) TestCreateGuest() {
	for i := 0; i < 4; i++ {
		var testUser models.User

		user, err := suite.dbService.CreateGuest()
		assert.Equal(suite.T(), nil, err)

		assert.NotEqual(suite.T(), "", user.Username)

		err = suite.dbService.DB.Where("username = ?", user.Username).First(&testUser).Error
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), user.Username, testUser.Username)
	}
}

type ChangePasswordTestData struct {
	username    string
	password    string
	newpassword string
	expectation string
}

var testChangePassword = []ChangePasswordTestData{
	{
		username:    "Admin",
		password:    "myPassword1",
		newpassword: "myOtherPassword1",
		expectation: "myOtherPassword1",
	},
	{
		username:    "Admin2",
		password:    "",
		newpassword: "myOtherPassword1",
		expectation: "myOtherPassword1",
	},
	{
		username:    "Admin3",
		password:    "myOtherPassword1",
		newpassword: "myOtherPassword1",
		expectation: "myOtherPassword1",
	},
}

func (suite *DatabaseTestSuite) TestChangePassword() {
	for _, test := range testChangePassword {
		err := suite.dbService.DB.Create(&models.User{
			Username: test.username,
		}).Error
		assert.Equal(suite.T(), nil, err)

		var user models.User
		err = suite.dbService.ChangePassword(test.username, test.password, test.newpassword)
		assert.Equal(suite.T(), nil, err)

		err = suite.dbService.DB.Where(models.User{
			Username: test.username,
		}).Select("username, password").First(&user).Error
		assert.Equal(suite.T(), nil, err)

		err = hash.CompareHashAndPassword([]byte(user.Password), []byte(test.expectation))
		assert.Equal(suite.T(), nil, err)

	}
	suite.dbService.DB.Close()
}

type RenameUsernameTestData struct {
	user        models.User
	newusername string
	expectation string
}

var testRenameUsername = []RenameUsernameTestData{
	{
		user: models.User{
			UserID:   uuid.UUID{},
			Username: "Admin1",
		},
		newusername: "AdminAdmin1",
		expectation: "AdminAdmin1",
	},
	{
		user: models.User{
			UserID:   uuid.UUID{},
			Username: "BE54-H474-PWB8-SD11",
		},
		newusername: "IwantAUsername",
		expectation: "IwantAUsername",
	},
}

func (suite *DatabaseTestSuite) TestRenameUser() {
	for _, test := range testRenameUsername {
		test.user.UserID = uuid.NewV4()
		err := suite.dbService.DB.Create(&test.user).Error
		assert.Equal(suite.T(), nil, err)

		err = suite.dbService.RenameUser(test.user.UserID, test.user.Username, test.newusername)
		assert.Equal(suite.T(), nil, err)

		var user models.User
		err = suite.dbService.DB.Where(models.User{
			UserID: test.user.UserID,
		}).Select("username").First(&user).Error
		assert.Equal(suite.T(), nil, err)

		assert.Equal(suite.T(), test.expectation, user.Username)
	}
	suite.dbService.DB.Close()
}

func (suite *DatabaseTestSuite) TestRenameUserNoNewUsername() {
	expectErr := "cannot change username to blank"
	err := suite.dbService.RenameUser(uuid.NewV4(), "Admin", "")
	assert.Equal(suite.T(), expectErr, err.Error())
	suite.dbService.DB.Close()
}

// func (suite *DatabaseTestSuite) TestRenameUserUsernameExists() {
// 	expectErr := "username alreay exists"
// 	id := uuid.NewV4()
// 	err := suite.dbService.DB.Create(&models.User{
// 		Username: "Admin1",
// 	}).Error
// 	assert.Equal(suite.T(), nil, err)
// 	err = suite.dbService.DB.Create(&models.User{
// 		UserID:   id,
// 		Username: "Admin2",
// 	}).Error
// 	assert.Equal(suite.T(), nil, err)
// 	fmt.Println(id)
// 	newErr := suite.dbService.RenameUser(id, "Admin2", "Admin1")
// 	assert.Equal(suite.T(), expectErr, newErr.Error())
// }
