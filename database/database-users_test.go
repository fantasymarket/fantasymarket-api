func (suite *DatabaseTestSuite) TestCreateGuest() {
	user, err := suite.dbService.CreateGuest()

	
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