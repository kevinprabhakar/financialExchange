package transaction

import (
	"financialExchange/sql"
	"financialExchange/util"
	"financialExchange/model"
	"errors"
	"time"
	"github.com/shopspring/decimal"
	gosql "database/sql"
	"fmt"
)

type TransactionController struct{
	Database 		*sql.MySqlDB
	Logger 			*util.Logger
	OrdersIncoming	chan model.OrderTransactionPackage
	ErrorChan 		chan error
}

func NewTransactionController(database *sql.MySqlDB, logger *util.Logger, ordersIncoming chan model.OrderTransactionPackage, errorChan chan error)(*TransactionController){
	return &TransactionController{
		Database: database,
		Logger: logger,
		OrdersIncoming: ordersIncoming,
		ErrorChan: errorChan,
	}
}

func (self *TransactionController) HandleMatchedOrders(){
	self.Logger.Debug("Spinning up worker to match orders")
	for orderTransactionPackage := range self.OrdersIncoming{
		//Grab Main Order
		mainOrder, mainOrderErr := self.Database.GetOrderById(orderTransactionPackage.MainOrder)
		if mainOrderErr != nil{
			self.ErrorChan <- errors.New("InvalidOrderId")
			self.Logger.ErrorMsg(mainOrderErr.Error())
			continue
		}

		//Grab Fulfilling Orders
		fulfillingOrders := make([]model.Order, 0)
		fulfillingOrdersIDs := make([]int64, 0)


		for key := range orderTransactionPackage.MatchingOrders{
			newOrder, err := self.Database.GetOrderById(key)
			if err != nil{
				self.ErrorChan <- errors.New("InvalidOrderId")
				self.Logger.ErrorMsg(err.Error())
				break
			}
			fulfillingOrders = append(fulfillingOrders, *newOrder)
			fulfillingOrdersIDs = append(fulfillingOrdersIDs, key)
		}


		//Update status, stockValue, cashValue, totalValue, and ownedShares of each order
		err := self.UpdateMainAndFulfillingOrders(*mainOrder, fulfillingOrders, orderTransactionPackage.MatchingOrders)
		if err != nil{
			self.ErrorChan <- errors.New("ErrorMovingMoney")
			self.Logger.ErrorMsg(err.Error())
			continue
		}

		//Create transaction
		transactionToBeInserted := model.Transaction{
			OrderPlaced: mainOrder.ID,
			OrdersFulfilling: fulfillingOrdersIDs,
			TotalCost: mainOrder.TotalCost,
			//System Fee is 0...FOR NOW
			SystemFee: model.NewMoneyObject(0.0),
			Created: time.Now().Unix(),
		}

		//Insert Transaction to Database
		_, err = self.Database.InsertTransactionToTable(transactionToBeInserted)
		if err != nil{
			self.ErrorChan <- err
			self.Logger.ErrorMsg(err.Error())
			continue
		}
		self.Logger.Debug(fmt.Sprintf("Matched order %d with %d orders", orderTransactionPackage.MainOrder, len(orderTransactionPackage.MatchingOrders)))
	}
}

func (self *TransactionController)UpdateMainAndFulfillingOrders(mainOrder model.Order,
																fulfillingOrders []model.Order,
																orderToShareMap map[int64]int)(error){
	numShares := 0

	tx, err := self.Database.Begin()
	if err != nil{
		return err
	}

	var newShareAmount int
	var newStockValue model.Money
	var newCashValue model.Money
	var newTotalValue model.Money
	updateStatusQuery := `UPDATE orders SET status = ?, fulfilled = ? WHERE id = ?`
	sqlQueryInsertOwnedShare := `INSERT INTO ownedShares (userId, security, numShares) VALUES ( ?, ?, ? )`


	for _, order := range fulfillingOrders{
		ownedShareForUser, err := self.Database.GetOwnedShareForUserForSecurity(order.Investor, order.Security)
		//Check to see if fulfillingOrder has an OwnedShare
		if err != nil{
			if err != gosql.ErrNoRows{
				return err
			}else{
				ownedShare := model.OwnedShare{
					UserID: order.Investor,
					Security: order.Security,
					NumShares: 0,
				}


				_, err := tx.Exec(sqlQueryInsertOwnedShare, ownedShare.UserID, ownedShare.Security, ownedShare.NumShares)
				if err != nil{
					tx.Rollback()
					return err
				}

				ownedShareForUser = ownedShare
			}
		}

		fulfillingPortfolio, err := self.Database.GetPortfolioByCustomerID(order.Investor)
		if err != nil{
			return err
		}

		//For each fulfilling order, update ownedShares, stockValue, and cashValue of user
		if mainOrder.InvestorAction == 0{
			//If Main Order is Buy, Fulfilling Orders are sells
			//Update Shares for Account
			newShareAmount = ownedShareForUser.NumShares-orderToShareMap[order.ID]
			//Subtract from Stock Value of account
			newStockValue = model.NewMoneyObjectFromDecimal(fulfillingPortfolio.StockValue.Sub((order.CostPerShare.Mul(decimal.NewFromFloat(float64(orderToShareMap[order.ID]))))))
			//Add to cash value of account
			newCashValue = model.NewMoneyObjectFromDecimal(fulfillingPortfolio.CashValue.Add((order.CostPerShare.Mul(decimal.NewFromFloat(float64(orderToShareMap[order.ID]))))))
			//Add two new values to get new total value for portfolio
			newTotalValue = model.NewMoneyObjectFromDecimal(newStockValue.Add(newCashValue.Decimal))


			//Create SQL Queries
			//Update OwnedShares
			updateOwnedSharesQuery := `UPDATE ownedShares SET numShares = ? WHERE security = ? AND userId = ?`
			_, err := tx.Exec(updateOwnedSharesQuery, newShareAmount, order.Security, order.Investor)
			if err != nil{
				tx.Rollback()
				return err
			}



			//Update Portfolio Amounts
			updatePortfolioAmounts := `UPDATE portfolios SET stockValue = ?, cashValue = ?, value = ? WHERE customer = ?`
			_, err = tx.Exec(updatePortfolioAmounts, newStockValue, newCashValue, newTotalValue, order.Investor)
			if err != nil{
				tx.Rollback()
				return err
			}




			//Update Order Book
			updateOrderBook := `UPDATE orders SET numSharesRemaining = ? WHERE id = ?`
			_, err = tx.Exec(updateOrderBook, newShareAmount, order.ID)
			if err != nil{
				tx.Rollback()
				return err
			}

			//Update status of order to Completed if newShareAmount == 0. else leave as is
			if newShareAmount == 0{

				_, err = tx.Exec(updateStatusQuery, 2, time.Now().Unix(), order.ID)
				if err != nil{
					tx.Rollback()
					return err
				}
			}

			//Increment numShares for mainOrder
			numShares += orderToShareMap[order.ID]
		}else{
			//Else if Main Order is Sell, Fulfilling Orders are buys
			//Update Shares for Account
			newShareAmount = ownedShareForUser.NumShares+orderToShareMap[order.ID]
			//Add to Stock Value of account

			newStockValue = model.NewMoneyObjectFromDecimal(fulfillingPortfolio.StockValue.Add((order.CostPerShare.Mul(decimal.NewFromFloat(float64(orderToShareMap[order.ID]))))))
			//Subtract from cash value of account
			newCashValue = model.NewMoneyObjectFromDecimal(fulfillingPortfolio.CashValue.Sub((order.CostPerShare.Mul(decimal.NewFromFloat(float64(orderToShareMap[order.ID]))))))
			//Add two new values to get new total value for portfolio
			newTotalValue = model.NewMoneyObjectFromDecimal(newStockValue.Add(newCashValue.Decimal))

			

			//Create SQL Queries
			//Update OwnedShares
			updateOwnedSharesQuery := `UPDATE ownedShares SET numShares = ? WHERE security = ? AND userId = ?`
			_, err := tx.Exec(updateOwnedSharesQuery, newShareAmount, order.Security, order.Investor)
			if err != nil{
				tx.Rollback()
				return err
			}

			


			//Update Portfolio Amounts
			updatePortfolioAmounts := `UPDATE portfolios SET stockValue = ?, cashValue = ?, value = ? WHERE customer = ?`
			_, err = tx.Exec(updatePortfolioAmounts, newStockValue, newCashValue, newTotalValue, order.Investor)
			if err != nil{
				tx.Rollback()
				return err
			}

			


			//Update Order Book
			updateOrderBook := `UPDATE orders SET numSharesRemaining = ? WHERE id = ?`
			_, err = tx.Exec(updateOrderBook, order.NumShares-newShareAmount, order.ID)
			if err != nil{
				tx.Rollback()
				return err
			}

			//If order is completed
			if newShareAmount == order.NumShares{
				

				_, err = tx.Exec(updateStatusQuery, 2, time.Now().Unix(), order.ID)
				if err != nil{
					tx.Rollback()
					return err
				}
			}

			//Increment numShares for mainOrder
			numShares += orderToShareMap[order.ID]
		}


	}


	//Grab current number of shares owned
	mainUserOwnedShares, mainUserErr := self.Database.GetOwnedShareForUserForSecurity(mainOrder.Investor, mainOrder.Security)

	if mainUserErr != nil{

		if mainUserErr != gosql.ErrNoRows{
			return mainUserErr
		}else{
			ownedShare := model.OwnedShare{
				UserID: mainOrder.Investor,
				Security: mainOrder.Security,
				NumShares: 0,
			}

			_, err := tx.Exec(sqlQueryInsertOwnedShare, ownedShare.UserID, ownedShare.Security, ownedShare.NumShares)
			if err != nil{
				tx.Rollback()
				return err
			}
		}
	}



	mainUserPortfolio, mainUserPortfolioErr := self.Database.GetPortfolioByCustomerID(mainOrder.Investor)

	if mainUserPortfolioErr != nil{
		return mainUserPortfolioErr
	}

	//Update ownedShares, stockValue and
	if mainOrder.InvestorAction == 0{
		//Increase num of ownedShares
		newShareAmount = mainUserOwnedShares.NumShares + numShares
		//Add to Stock Value of account
		newStockValue = model.NewMoneyObjectFromDecimal(mainUserPortfolio.StockValue.Add((mainOrder.CostPerShare.Mul(decimal.NewFromFloat(float64(numShares))))))
		//Subtract from cash value of account
		newCashValue = model.NewMoneyObjectFromDecimal(mainUserPortfolio.CashValue.Sub((mainOrder.CostPerShare.Mul(decimal.NewFromFloat(float64(numShares))))))
		//Add stock and cash value to get new total value
		newTotalValue = model.NewMoneyObjectFromDecimal(newCashValue.Add(newStockValue.Decimal))


		//Create SQL Queries
		//Update OwnedShares
		updateOwnedSharesQuery := `UPDATE ownedShares SET numShares = ? WHERE security = ? AND userId = ?`
		_, err := tx.Exec(updateOwnedSharesQuery, newShareAmount, mainOrder.Security, mainOrder.Investor)
		if err != nil{
			tx.Rollback()
			return err
		}




		//Update Portfolio Amounts
		updatePortfolioAmounts := `UPDATE portfolios SET stockValue = ?, cashValue = ?, value = ? WHERE customer = ?`
		_, err = tx.Exec(updatePortfolioAmounts, newStockValue, newCashValue, newTotalValue, mainOrder.Investor)
		if err != nil{
			tx.Rollback()
			return err
		}



		//Update Order Book
		updateOrderBook := `UPDATE orders SET numSharesRemaining = ? WHERE id = ?`
		_, err = tx.Exec(updateOrderBook, mainOrder.NumShares-newShareAmount, mainOrder.ID)
		if err != nil{
			tx.Rollback()
			return err
		}

		//Update mainOrder status to Complete
		if numShares == mainOrder.NumSharesRemaining{

			_, err = tx.Exec(updateStatusQuery, 2, time.Now().Unix(), mainOrder.ID)
			if err != nil{
				tx.Rollback()
				return err
			}
		}

	}else{
		//Decrease num of ownedShares
		newShareAmount = mainUserOwnedShares.NumShares - numShares
		//Subtract from Stock Value of account
		newStockValue = model.NewMoneyObjectFromDecimal(mainUserPortfolio.StockValue.Sub((mainOrder.CostPerShare.Mul(decimal.NewFromFloat(float64(numShares))))))
		//Add to cash value of account
		newCashValue = model.NewMoneyObjectFromDecimal(mainUserPortfolio.CashValue.Add((mainOrder.CostPerShare.Mul(decimal.NewFromFloat(float64(numShares))))))
		//Add stock and cash value to get new total value
		newTotalValue = model.NewMoneyObjectFromDecimal(newCashValue.Add(newStockValue.Decimal))

		//Create SQL Queries
		//Update OwnedShares
		updateOwnedSharesQuery := `UPDATE ownedShares SET numShares = ? WHERE security = ? AND userId = ?`
		_, err := tx.Exec(updateOwnedSharesQuery, newShareAmount, mainOrder.Security, mainOrder.Investor)
		if err != nil{

			tx.Rollback()
			return err
		}

		//Update Portfolio Amounts
		updatePortfolioAmounts := `UPDATE portfolios SET stockValue = ?, cashValue = ?, value = ? WHERE customer = ?`
		_, err = tx.Exec(updatePortfolioAmounts, newStockValue, newCashValue, newTotalValue, mainOrder.Investor)
		if err != nil{
			tx.Rollback()
			return err
		}

		//Update Order Book
		updateOrderBook := `UPDATE orders SET numSharesRemaining = ? WHERE id = ?`
		_, err = tx.Exec(updateOrderBook, newShareAmount, mainOrder.ID)
		if err != nil{
			tx.Rollback()
			return err
		}
		//Update mainOrder status to Complete
		if newShareAmount == 0{
			_, err = tx.Exec(updateStatusQuery, 2, time.Now().Unix(),mainOrder.ID)
			if err != nil{
				tx.Rollback()
				return err
			}
		}
	}



	//Commit SQL Transaction to database
	err = tx.Commit()
	if err != nil{

		tx.Rollback()
		return err
	}

	return nil
}