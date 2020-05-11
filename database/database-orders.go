package database

import (
	"errors"
	"fantasymarket/database/models"
	"fantasymarket/utils"
	"time"

	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

var (
	// ErrOrderFilledOrCancelled means a order is already filled or cancelled
	ErrOrderFilledOrCancelled = errors.New("can't cancel order, as its already filled or cancelled")
	// ErrInvalidAmount means the order's amount is invalid
	ErrInvalidAmount = errors.New("amount cannot be less than 0")
	// ErrNotEnoughMoney means the order can't be executed
	ErrNotEnoughMoney = errors.New("insufficient balance")
	// ErrInvalidType means the type given in the order is not found
	ErrInvalidType = errors.New("404 type not found")
	// ErrCantSellMoreThanYouHave meansyou can't sell more than
	ErrCantSellMoreThanYouHave = errors.New("cant sell more than you have")

	// ListOFValidTypes is the list of types accepted for trading
	ListOFValidTypes = [3]string{"stock", "crypto", "earth"}
)

// AddOrder adds an Order to the database
func (s *Service) AddOrder(order models.Order, userID uuid.UUID, currentDate time.Time) error {
	return s.DB.Create(&models.Order{
		UserID:    userID,
		CreatedAt: currentDate,
		Type:      order.Type,
		Side:      order.Side,
		Symbol:    order.Symbol,
		Status:    order.Status,
	}).Error
}

// GetOrders gets all orders based on the parameters of orderDetails where Symbol, Type and userID can be set.
// Limit is how many items. Offset is from where to where the data is used
func (s *Service) GetOrders(orderDetails models.Order, limit int, offset int) (*[]models.Order, error) {
	var orders *[]models.Order
	if err := s.DB.Where(models.Order{UserID: orderDetails.UserID, Type: orderDetails.Type, Symbol: orderDetails.Symbol}).Limit(limit).Offset(offset).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID gets an order by the orderID
func (s *Service) GetOrderByID(orderID uuid.UUID) (*models.Order, error) {
	var order *models.Order
	if err := s.DB.Where(models.Order{OrderID: orderID}).First(&order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

// CancelOrder cancels an order in the database
func (s *Service) CancelOrder(orderID uuid.UUID, currentDate time.Time) error {

	var order models.Order
	if err := s.DB.Where(models.Order{OrderID: orderID}).First(&order).Error; err != nil {
		return err
	}

	// TODO check if the order is still active
	if order.Status == "waiting" {
		return s.DB.Model(&order).Updates(models.Order{Status: "cancelled", FilledAt: currentDate}).Error
	}

	return ErrOrderFilledOrCancelled
}

// DONT READ BEYOND THIS LINE
// HERE BE DRAGONS

//  _   _     _       _        __ _
// | | | |   (_)     (_)      / _(_)
// | |_| |__  _ ___   _ ___  | |_ _ _ __   ___
// | __| '_ \| / __| | / __| |  _| | '_ \ / _ \
// | |_| | | | \__ \ | \__ \ | | | | | | |  __/
//  \__|_| |_|_|___/ |_|___/ |_| |_|_| |_|\___|

// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
// @@@@@@@@@@@@@@@@@@#####@@@@@@@@@@@@@#@@@@@@@@@@@@@@@###@@@@@@@@@@@@@@@@@#######@
// ###@@@@@@@@@@#############################@@@@@@###########@@@@@@@@@############
// ###@#######@###############@@##########@################@@@@@@@###@####@#####@@@
// ############@###@@@@@@@#####@..+@@@+`    @@########@@@#####@##@@###@@@##########
// #######@#####@@@#######@@@;                      ;;;;;;';'';;;;@@@######@+  ;@@@
// @####@:                                          '';#########+;;    ``
//                              `,:''               ;;;#########+;@
//             @@@@@@@@@@@@@@@@@@@@@@               ;;;##@######';#
//             @                  @@@               ';;####@####;;@
//       @`    @        `      `  @@@               ;;;######@##''@
//      @@@    @     #@@       @;::@@               ;;;#########;'@          @
//     @@@@@   @   #@@@.   @  @:::;+@@:'`           ;;;#####@##@;;@          @@:
//    @@@@@@   @  @@@@@    @@  ;@@@@;::;`           ;';''''''+#;;;.          .@@`
//   @@@@@@#   @ @@@@@@   @@@@::::::::#+;;;@        #;;;;;;;;;;;;@            @@@
//   @@@@@@@   @@@@,@@@  ;@@@@@:::;+@;;::::@'#.                              ,@@@
// @.@@@,@@@   @@@,,:@@@ @@:,@@ @;;:+#;;:+  @@@@                   '         @@@@`@
// @@@@@,@@@ ' @@@,,,@@@@@@,,,@@@;+ .. @;@ @@@@+                 @@         @@@@@@@
// @@@@@,@@  @+'@,,,,,@@@@,,@@@@@: '@@@ @+` @@@@#@@@@@@+       ;@@#        @@@:@@@@
// @@@@@,@@  @@#,,,,,,,:@,,@@@@@@: @@@@ #:;+::::;@@@@@@@      @@#@@       @@@,,@@@@
// ,@@@@,@@  @@@:,,,,,,,,,@@@@@@@;+   `@;;:;::::;;@@@@@     `@@,'@@+     #@@#,,+@@,
// ,,@@,,@@ ,@@@@,,,,,,,,:@@@@@@:::::::::;+##:;;;;::++ ,   @@@,@++++++#@ @@@:,,,@,,
// ,,,,,@@+ @@@@@,,,,,,,+''@@@#;;:::::;:;:::;@  ..`   @@  @@@,,@,;@@@+,@.,@@,,,,,,,
// ,,,,,@@@ @@@@@,,,,,,:+''@'@;::;;:;;:::::::@       @@@ @#@,,,#,......@,,
// ,,,,,@@@#@@@@@,,,,,,,'''@'@;;;@;@;::::::;;'@     .@@@@#@,,  @,......@'
// ,,,,,'@@@@@,@@,,,,,::+''@'@;;;+;@;:::::::;;'     '@@:@@;    @:....,.@
// ,,,,,,@@@@@,@@::::::,@''@''';;:+;#:::;:::;+:,    :@@,,@,,
// ,,,,,,:@@@,,@@:::::::@''''+@';;;#;'@@@@@'#:#@#   `@@:,,,,,,,,,,.,,,,,,,,,,,,,,,,
// ,,,,,,,+@,,,@@;;;;'',:@''''''@#;;;;;#@;;@''#:; , `@@+,,,,,,,,,,,,,,,,,,,,,,,,,,,
// ,,,,@,,,:,,:@@::::,,::::@''''''''''''@;::#@@;:;;@.@@#,,,,,,,,,,,.,,@@,,,,,,,,,,@
// ,,,,;,,,,,,:@@::::@@,:::,@@@####@@@@@@;:;;#@#@@  @@@,,,,,,,,@@@@##@@@########@@@
// ,,,,:,,,,,,,@@'::@@@:::::::#''+++++@@@##@@@@@@@#@@@@,,,@@,,,@@@@@@@@@#######@@@@
// ,,,,,,,,,,,,@@@,@@@@@:,::::+''+::,,@@@'':+@@@:::@@,,;@:@,,,,@@@@@@@@@######@@@@:
// ,,;,,+,,,,,,,@@@@@'@@,:,@,,,''+::::@@@'':'@@@::,:::::,@#+:,#@@@@@@@@#####@@@@@+,
// ,,+,,@,,#,@,,,@@@@,@@@,@@:::''+:::,@@@''::@@@::::::::::::@@@#@+,@@@@####@@@@@@,,
// ,,,,,;,,:,,,,,,@@@,,@@@@@:::'''::::@@@'',:@@@:::::::,::#@@@@@,,,@@@@##@@@@@@@,,,
// ,,,,,,',,#,;,,,,:@,,,#@@@:::#@:::::,:@''::@@,::::::::;@@@@@;,,,,@@@@@@@@@@@:,,,,
// ,,,,,,,,,,,,,,,,,,,,,,,@##:,,::::::::,::::::::::::::,@@@#@:,,,,,:@@@@@@@@,,,,,,,
// ,,,,,,,,,,,,,,,,,,,,,,,'@:::::::::::,::::::::::::::,@@@@@@,,,,,,,,'@@@@,,,,,,,,,
// ,,,,,,,,,,,,,,,,,,,,,:@::::::::::::::::::::::::::::,@@@@@,,,,,,,,,,,;#,,,,,,,,,,

// FillOrder "completes" the order
// Note: this is only supposed to be called from the game code
// This should ONLY be called if the order is actually supposed to run
// and things like the order type, limit prices etc. have been checked
// (since we only check if the user has enough money in his portfolio)
func (s *Service) FillOrder(orderID uuid.UUID, userID uuid.UUID, currentIndex int64, currentDate time.Time) error {
	order, err := s.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	var user models.User
	if err := s.DB.Where(models.User{
		UserID: userID,
	}).Preload("Portfolio").Find(&user).Error; err != nil {
		return err
	}

	if order.Amount < 0 {
		return ErrInvalidAmount
	}

	if !utils.Includes(ListOFValidTypes[:], order.Type) {
		return ErrInvalidType
	}

	if order.Side == "sell" {
		order.Amount = order.Amount * -1
	}

	price := order.Amount * currentIndex
	if (user.Portfolio.Balance + price) < 0 {
		s.CancelOrder(order.OrderID, currentDate)
		return ErrNotEnoughMoney
	}

	var affectedPortfolioItem models.PortfolioItem
	if err := s.DB.Where(models.PortfolioItem{
		PortfolioID: user.Portfolio.PortfolioID,
		Symbol:      order.Symbol,
	}).Attrs(models.PortfolioItem{
		PortfolioID: user.Portfolio.PortfolioID,
		Symbol:      order.Symbol,
		Type:        order.Type,
		Amount:      0,
	}).FirstOrCreate(&affectedPortfolioItem).Error; err != nil {
		s.CancelOrder(order.OrderID, currentDate)
		return err
	}

	newAmount := order.Amount + affectedPortfolioItem.Amount
	if newAmount < 0 {
		s.CancelOrder(order.OrderID, currentDate)
		return ErrCantSellMoreThanYouHave
	}

	newBalance := user.Portfolio.Balance - price
	newPortfolio, err := s.updatePortfolioItem(
		user.Portfolio.PortfolioID,
		affectedPortfolioItem.PortfolioItemID,
		newAmount,
		newBalance,
	)

	if err != nil || s.createPortfolioSnapshot(newPortfolio, currentDate) != nil {
		s.CancelOrder(order.OrderID, currentDate)
		return err
	}

	return s.DB.Where(models.Order{OrderID: orderID}).Updates(models.Order{Status: "filled", FilledAt: currentDate}).Error
}

func (s *Service) createPortfolioSnapshot(portfolio *models.Portfolio, currentDate time.Time) error {
	snapshotData, err := json.Marshal(portfolio)
	if err != nil {
		return err
	}

	return s.DB.Create(&models.PortfolioSnapshot{
		UserID:      portfolio.UserID,
		PortfolioID: portfolio.PortfolioID,
		CreatedAt:   currentDate,
		Data:        string(snapshotData),
	}).Error
}

func (s *Service) updatePortfolioItem(portfolioID uuid.UUID, itemID uuid.UUID, newAmount int64, newBalance int64) (*models.Portfolio, error) {
	// we use a transactions since if updating
	// balance fails, we also need to rollback
	// the portfolioItem's amount
	tx := s.DB.Begin()

	if err := tx.Where(&models.PortfolioItem{
		PortfolioItemID: itemID,
	}).Updates(models.PortfolioItem{
		Amount: newAmount,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var portfolio models.Portfolio
	if err := tx.Where(&models.Portfolio{
		PortfolioID: portfolioID,
	}).Find(&portfolio).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(portfolio).Updates(models.Portfolio{
		Balance: newBalance,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &portfolio, tx.Commit().Error
}
