package testing

import (
	"testing"
	"financialExchange/model"
	"time"
	"math/rand"
	"fmt"
)

//func TestDropTables(t *testing.T){
//	err := TestingController.ResetTables()
//	if err != nil{
//		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
//		t.Fail()
//		return
//	}
//}

func TestOneBuyOneSell(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	_, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	buyUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Buy User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserMoney(buyUserID, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving buy user money")
		t.Fail()
		return
	}

	sellUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Sell User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserStock(sellUserID,securityID, 1, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving sell user stock")
		t.Fail()
		return
	}

	buyOrderID, err := TestingController.PlaceBuyOrderForUser(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for buy user")
		t.Fail()
		return
	}

	sellOrderID, err := TestingController.PlaceSellOrderForUser(sellUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for sell user")
		t.Fail()
		return
	}

	//Momentary sleep for order matching
	time.Sleep(300 * time.Millisecond)

	//Validating Database
	buyerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(buyUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buyer portfolio")
		t.Fail()
		return
	}

	sellerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(sellUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing seller portfolio")
		t.Fail()
		return
	}

	buyerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(buyUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for buy user")

		t.Fail()
		return
	}

	sellerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(sellUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for sell user")
		t.Fail()
		return
	}
	buyOrder, err := TestingController.Db.GetOrderById(buyOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buy order")
		t.Fail()
		return
	}

	sellOrder, err := TestingController.Db.GetOrderById(sellOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing sell order")
		t.Fail()
		return
	}

	if !(buyerPortfolio.StockValue.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.CashValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Buyer has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(sellerPortfolio.CashValue.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.StockValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Seller has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(buyerOwnedShares.NumShares == 1){
		TestingController.Logger.ErrorMsg("Buyer has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(sellerOwnedShares.NumShares == 0){
		TestingController.Logger.ErrorMsg("Seller has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(buyOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- buy order has incorrect completion status")
		t.Fail()
		return
	}
	if !(sellOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- sell order has incorrect completion status")
		t.Fail()
		return
	}
}

func TestOneSellOneBuy(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	_, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	buyUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Buy User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserMoney(buyUserID, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving buy user money")
		t.Fail()
		return
	}

	sellUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Sell User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserStock(sellUserID,securityID, 1, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving sell user stock")
		t.Fail()
		return
	}

	sellOrderID, err := TestingController.PlaceSellOrderForUser(sellUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for sell user")
		t.Fail()
		return
	}

	buyOrderID, err := TestingController.PlaceBuyOrderForUser(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for buy user")
		t.Fail()
		return
	}

	//Momentary sleep for order matching
	time.Sleep(300 * time.Millisecond)

	//Validating Database
	buyerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(buyUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buyer portfolio")
		t.Fail()
		return
	}

	sellerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(sellUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing seller portfolio")
		t.Fail()
		return
	}

	buyerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(buyUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for buy user")

		t.Fail()
		return
	}

	sellerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(sellUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for sell user")
		t.Fail()
		return
	}
	buyOrder, err := TestingController.Db.GetOrderById(buyOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buy order")
		t.Fail()
		return
	}

	sellOrder, err := TestingController.Db.GetOrderById(sellOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing sell order")
		t.Fail()
		return
	}

	if !(buyerPortfolio.StockValue.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.CashValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Buyer has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(sellerPortfolio.CashValue.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.StockValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Seller has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(buyerOwnedShares.NumShares == 1){
		TestingController.Logger.ErrorMsg("Buyer has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(sellerOwnedShares.NumShares == 0){
		TestingController.Logger.ErrorMsg("Seller has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(buyOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- buy order has incorrect completion status")
		t.Fail()
		return
	}
	if !(sellOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- sell order has incorrect completion status")
		t.Fail()
		return
	}
}

func TestTenBuyOneSell(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	_, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	buyUserIds := make([]int64, 0)
	buyOrderIds := make([]int64, 0)

	//Create 10 Buy Users, give them each $100, and execute a buy order for all of them
	for i:=0; i < 10; i++{
		buyUserID, err := TestingController.CreateUser()
		if err != nil{
			TestingController.Logger.ErrorMsg("Error Creating Buy Users")
			t.Fail()
			return
		}
		err = TestingController.GiveUserMoney(buyUserID, model.NewMoneyObject(100.00))
		if err != nil{
			TestingController.Logger.ErrorMsg("Error giving buy user money")
			t.Fail()
			return
		}
		buyOrderID, err := TestingController.PlaceBuyOrderForUser(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
		if err != nil{
			TestingController.Logger.ErrorMsg("Error placing order for buy user")
			t.Fail()
			return
		}
		buyUserIds = append(buyUserIds, buyUserID)
		buyOrderIds = append(buyOrderIds, buyOrderID)
	}

	sellUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Sell User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserStock(sellUserID,securityID, 10, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving sell user stock")
		t.Fail()
		return
	}

	sellOrderID, err := TestingController.PlaceSellOrderForUser(sellUserID, securityID, symbol, model.NewMoneyObject(100.00),10)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for sell user")
		t.Fail()
		return
	}

	//Momentary sleep for order matching
	time.Sleep(300 * time.Millisecond)

	//Validating Database
	for index := 0; index < len(buyOrderIds); index++{
		buyerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(buyUserIds[index])
		if err != nil{
			TestingController.Logger.ErrorMsg("Error grabbing buyer portfolio")
			t.Fail()
			return
		}

		buyerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(buyUserIds[index], securityID)
		if err != nil{
			TestingController.Logger.ErrorMsg("Error grabbing owned shares for buy user")

			t.Fail()
			return
		}
		buyOrder, err := TestingController.Db.GetOrderById(buyOrderIds[index])
		if err != nil{
			TestingController.Logger.ErrorMsg("Error grabbing buy order")
			t.Fail()
			return
		}
		if !(buyerPortfolio.StockValue.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.CashValue.Equals(model.NewMoneyObject(0.0).Decimal)){
			TestingController.Logger.ErrorMsg("Buyer has incorrect portfolio values")
			t.Fail()
			return
		}
		if !(buyerOwnedShares.NumShares == 1){
			TestingController.Logger.ErrorMsg("Buyer has incorrect ownedShare values")
			t.Fail()
			return
		}
		if !(buyOrder.Status == 2){
			TestingController.Logger.ErrorMsg("Error -- buy order has incorrect completion status")
			t.Fail()
			return
		}
	}

	sellerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(sellUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing seller portfolio")
		t.Fail()
		return
	}

	sellerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(sellUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for sell user")
		t.Fail()
		return
	}


	sellOrder, err := TestingController.Db.GetOrderById(sellOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing sell order")
		t.Fail()
		return
	}

	if !(sellerPortfolio.CashValue.Equals(model.NewMoneyObject(1000.00).Decimal) && sellerPortfolio.Value.Equals(model.NewMoneyObject(1000.00).Decimal) && sellerPortfolio.StockValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Seller has incorrect portfolio values")
		t.Fail()
		return
	}

	if !(sellerOwnedShares.NumShares == 0){
		TestingController.Logger.ErrorMsg("Seller has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(sellOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- sell order has incorrect completion status")
		t.Fail()
		return
	}
}

func TestTenSellOneBuy(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	_, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	sellUserIds := make([]int64, 0)
	sellOrderIds := make([]int64, 0)

	//Create 10 sell Users, give 1 share at $100/share, and execute a sell order for all of them
	for i:=0; i < 10; i++{
		sellUserID, err := TestingController.CreateUser()
		if err != nil{
			TestingController.Logger.ErrorMsg("Error Creating sell Users")
			t.Fail()
			return
		}
		err = TestingController.GiveUserStock(sellUserID,securityID, 1, model.NewMoneyObject(100.00))
		if err != nil{
			TestingController.Logger.ErrorMsg("Error giving sell user stock")
			t.Fail()
			return
		}
		sellOrderID, err := TestingController.PlaceSellOrderForUser(sellUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
		if err != nil{
			TestingController.Logger.ErrorMsg("Error placing order for sell user")
			t.Fail()
			return
		}
		sellUserIds = append(sellUserIds, sellUserID)
		sellOrderIds = append(sellOrderIds, sellOrderID)
	}

	buyUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating buy User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserMoney(buyUserID,model.NewMoneyObject(1000.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving buy user money")
		t.Fail()
		return
	}

	sellOrderID, err := TestingController.PlaceBuyOrderForUser(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),10)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for buy user")
		t.Fail()
		return
	}

	//Momentary sleep for order matching
	time.Sleep(300 * time.Millisecond)

	//Validating Database
	for index := 0; index < len(sellOrderIds); index++{
		sellerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(sellUserIds[index])
		if err != nil{
			TestingController.Logger.ErrorMsg("Error grabbing seller portfolio")
			t.Fail()
			return
		}

		sellerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(sellUserIds[index], securityID)
		if err != nil{
			TestingController.Logger.ErrorMsg("Error grabbing owned shares for sell user")

			t.Fail()
			return
		}
		sellOrder, err := TestingController.Db.GetOrderById(sellUserIds[index])
		if err != nil{
			TestingController.Logger.ErrorMsg("Error grabbing sell order")
			t.Fail()
			return
		}
		if !(sellerPortfolio.StockValue.Equals(model.NewMoneyObject(0.0).Decimal) && sellerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.CashValue.Equals(model.NewMoneyObject(100.0).Decimal)){
			TestingController.Logger.ErrorMsg("Seller has incorrect portfolio values")
			t.Fail()
			return
		}
		if !(sellerOwnedShares.NumShares == 0){
			TestingController.Logger.ErrorMsg("Seller has incorrect ownedShare values")
			t.Fail()
			return
		}
		if !(sellOrder.Status == 2){
			TestingController.Logger.ErrorMsg("Error -- sell order has incorrect completion status")
			t.Fail()
			return
		}
	}

	buyerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(buyUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buyer portfolio")
		t.Fail()
		return
	}

	buyerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(buyUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for buy user")
		t.Fail()
		return
	}


	buyerOrder, err := TestingController.Db.GetOrderById(sellOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buy order")
		t.Fail()
		return
	}

	if !(buyerPortfolio.CashValue.Equals(model.NewMoneyObject(0.00).Decimal) && buyerPortfolio.Value.Equals(model.NewMoneyObject(1000.00).Decimal) && buyerPortfolio.StockValue.Equals(model.NewMoneyObject(1000.0).Decimal)){
		TestingController.Logger.ErrorMsg("Buyer has incorrect portfolio values")
		t.Fail()
		return
	}

	if !(buyerOwnedShares.NumShares == 10){
		TestingController.Logger.ErrorMsg("Buyer has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(buyerOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- buy order has incorrect completion status")
		t.Fail()
		return
	}
}

func TestOneBuyOneIncompleteSell(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	_, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	buyUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Buy User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserMoney(buyUserID, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving buy user money")
		t.Fail()
		return
	}

	sellUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Sell User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserStock(sellUserID,securityID, 3, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving sell user stock")
		t.Fail()
		return
	}

	buyOrderID, err := TestingController.PlaceBuyOrderForUser(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for buy user")
		t.Fail()
		return
	}

	sellOrderID, err := TestingController.PlaceSellOrderForUser(sellUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for sell user")
		t.Fail()
		return
	}

	//Momentary sleep for order matching
	time.Sleep(300 * time.Millisecond)

	//Validating Database
	buyerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(buyUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buyer portfolio")
		t.Fail()
		return
	}

	sellerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(sellUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing seller portfolio")
		t.Fail()
		return
	}

	buyerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(buyUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for buy user")

		t.Fail()
		return
	}

	sellerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(sellUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for sell user")
		t.Fail()
		return
	}
	buyOrder, err := TestingController.Db.GetOrderById(buyOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buy order")
		t.Fail()
		return
	}

	sellOrder, err := TestingController.Db.GetOrderById(sellOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing sell order")
		t.Fail()
		return
	}

	if !(buyerPortfolio.StockValue.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.CashValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Buyer has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(sellerPortfolio.CashValue.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.Value.Equals(model.NewMoneyObject(300.00).Decimal) && sellerPortfolio.StockValue.Equals(model.NewMoneyObject(200.0).Decimal)){
		TestingController.Logger.ErrorMsg("Seller has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(buyerOwnedShares.NumShares == 1){
		TestingController.Logger.ErrorMsg("Buyer has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(sellerOwnedShares.NumShares == 2){
		TestingController.Logger.ErrorMsg("Seller has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(buyOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- buy order has incorrect completion status")
		t.Fail()
		return
	}
	if !(sellOrder.Status == 1){
		TestingController.Logger.ErrorMsg("Error -- sell order has incorrect completion status")
		t.Fail()
		return
	}
}

func TestOneIncompleteBuyOneSell(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	_, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	buyUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Buy User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserMoney(buyUserID, model.NewMoneyObject(300.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving buy user money")
		t.Fail()
		return
	}

	sellUserID, err := TestingController.CreateUser()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Sell User")
		t.Fail()
		return
	}

	err = TestingController.GiveUserStock(sellUserID,securityID, 1, model.NewMoneyObject(100.00))
	if err != nil{
		TestingController.Logger.ErrorMsg("Error giving sell user stock")
		t.Fail()
		return
	}

	buyOrderID, err := TestingController.PlaceBuyOrderForUser(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),3)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for buy user")
		t.Fail()
		return
	}

	sellOrderID, err := TestingController.PlaceSellOrderForUser(sellUserID, securityID, symbol, model.NewMoneyObject(100.00),1)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error placing order for sell user")
		t.Fail()
		return
	}

	//Momentary sleep for order matching
	time.Sleep(300 * time.Millisecond)

	//Validating Database
	buyerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(buyUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buyer portfolio")
		t.Fail()
		return
	}

	sellerPortfolio, err := TestingController.Db.GetPortfolioByCustomerID(sellUserID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing seller portfolio")
		t.Fail()
		return
	}

	buyerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(buyUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for buy user")

		t.Fail()
		return
	}

	sellerOwnedShares, err := TestingController.Db.GetOwnedShareForUserForSecurity(sellUserID, securityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing owned shares for sell user")
		t.Fail()
		return
	}
	buyOrder, err := TestingController.Db.GetOrderById(buyOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing buy order")
		t.Fail()
		return
	}

	sellOrder, err := TestingController.Db.GetOrderById(sellOrderID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error grabbing sell order")
		t.Fail()
		return
	}

	if !(buyerPortfolio.StockValue.Equals(model.NewMoneyObject(100.00).Decimal) && buyerPortfolio.Value.Equals(model.NewMoneyObject(300.00).Decimal) && buyerPortfolio.CashValue.Equals(model.NewMoneyObject(200.0).Decimal)){
		TestingController.Logger.ErrorMsg("Buyer has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(sellerPortfolio.CashValue.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.Value.Equals(model.NewMoneyObject(100.00).Decimal) && sellerPortfolio.StockValue.Equals(model.NewMoneyObject(0.0).Decimal)){
		TestingController.Logger.ErrorMsg("Seller has incorrect portfolio values")
		t.Fail()
		return
	}
	if !(buyerOwnedShares.NumShares == 1){
		TestingController.Logger.ErrorMsg("Buyer has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(sellerOwnedShares.NumShares == 0){
		TestingController.Logger.ErrorMsg("Seller has incorrect ownedShare values")
		t.Fail()
		return
	}
	if !(buyOrder.Status == 1){
		TestingController.Logger.ErrorMsg("Error -- buy order has incorrect completion status")
		t.Fail()
		return
	}
	if !(sellOrder.Status == 2){
		TestingController.Logger.ErrorMsg("Error -- sell order has incorrect completion status")
		t.Fail()
		return
	}

}

func TestRandomWithIPO(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}

	entityID, securityID, symbol, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	_, err = TestingController.IPOEntity(entityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Doing IPO")
		t.Fail()
		return
	}

	buyUserIds := make([]int64, 0)
	buyUserOrders := make([]int64, 0)

	for i:=0; i < 10; i++{
		buyUserID, err := TestingController.CreateUser()
		if err != nil{
			TestingController.Logger.ErrorMsg("Error Creating Buy Users")
			t.Fail()
			return
		}
		err = TestingController.GiveUserMoney(buyUserID, model.NewMoneyObject(100.00))
		if err != nil{
			TestingController.Logger.ErrorMsg("Error giving buy user money")
			t.Fail()
			return
		}
		//86400
		s2 := rand.NewSource(time.Now().UnixNano())
		r2 := rand.New(s2)
		randNum := r2.Int63n(86400)-86400
		fmt.Println(randNum)
		buyOrderID, err := TestingController.PlaceBuyOrderForUserWithSpecTime(buyUserID, securityID, symbol, model.NewMoneyObject(100.00),1,time.Now().Unix()+randNum)
		if err != nil{
			TestingController.Logger.ErrorMsg("Error placing order for buy user")
			t.Fail()
			return
		}
		buyUserIds = append(buyUserIds, buyUserID)
		buyUserOrders = append(buyUserOrders, buyOrderID)
	}
}

func TestBrownianForFrontEnd(t *testing.T){
	err := TestingController.ResetTables()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Resetting SQL tables")
		t.Fail()
		return
	}
	entityID, _, _, err := TestingController.CreateEntityWithSecurity()
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Creating Entity")
		t.Fail()
		return
	}

	_, err = TestingController.IPOEntity(entityID)
	if err != nil{
		TestingController.Logger.ErrorMsg("Error Doing IPO")
		t.Fail()
		return
	}

	//order, err := TestingController.Db.GetOrderById(int64(1))
	//if err != nil{
	//	TestingController.Logger.ErrorMsg("Error Getting Entity")
	//	t.Fail()
	//	return
	//}
	//
	//initialPrice, _ := order.CostPerShare.Float64()
	//
	//artificialPrices, _, err := TestingController.GeometricBrownianMotion(initialPrice,0.1,0.4,1.0,288)
	//if err != nil{
	//	TestingController.Logger.ErrorMsg("Error Generating Brownian Motion")
	//	t.Fail()
	//	return
	//}
	//
	//err = TestingController.CreateTransactionsForEntityFromGBM(entityID, artificialPrices)
	//if err != nil{
	//	TestingController.Logger.ErrorMsg("Error Generating Brownian Motion Transactions")
	//	t.Fail()
	//	return
	//}





}
