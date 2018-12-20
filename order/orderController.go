package order

import (
	"financialExchange/util"
	"financialExchange/model"
	"financialExchange/sql"
	"errors"
	gosql "database/sql"
	"golang.org/x/net/html"

	"time"

	"fmt"
)

type OrderController struct {
	Database        *sql.MySqlDB
	Logger         *util.Logger
	OrdersOutgoing chan model.OrderTransactionPackage
}

func NewOrderController (database *sql.MySqlDB, logger *util.Logger, ordersOutgoing chan model.OrderTransactionPackage)(*OrderController){
	return &OrderController{
		Database: database,
		Logger: logger,
		OrdersOutgoing: ordersOutgoing,
	}
}

func (self *OrderController) CreateOrder(params model.OrderCreateParams)(int64, error){
	//Validate order fields
	validErr := self.ValidateOrderFields(params)

	if validErr != nil{
		return 0, validErr
	}

	//Get security by symbol
	securityId, err := self.Database.CheckIfSecurityInTable(params.Symbol)
	if err != nil{
		return 0, errors.New("Invalid Symbol")
	}

	//Create Order
	orderToInsert := model.Order{
		Investor: params.UserID,
		InvestorAction: model.InvestorAction(params.InvestorAction),
		InvestorType: model.InvestorType(params.InvestorType),
		OrderType: model.OrderType(params.OrderType),
		Security: securityId,
		Symbol: params.Symbol,
		NumShares: params.NumShares,
		NumSharesRemaining: params.NumShares,
		CostPerShare: model.NewMoneyObject(params.CostPerShare),
		CostOfShares: model.NewMoneyObject(float64(params.NumShares) * params.CostPerShare),
		//No system fee... FOR NOW
		SystemFee: model.NewMoneyObject(0.0),
		//Total Cost = (Cost of Shares * Cost Per Share) + 0.0
		TotalCost: model.NewMoneyObject((float64(params.NumShares) * params.CostPerShare)+0.0),
		Created: params.TimeCreated,
		Updated: params.TimeCreated,
		Fulfilled: params.TimeCreated,
		Status: model.Untouched,
		AllowTakers: params.AllowTakers,
		LimitPerShare: model.NewMoneyObject(params.LimitPerShare),
		TakerFee: model.NewMoneyObject(0.0),
		StopPrice: model.NewMoneyObject(0.0),
	}

	//Insert order into database
	insertId, insertErr := self.Database.InsertOrderIntoTable(orderToInsert)

	if insertErr != nil{
		return 0, insertErr
	}

	//Get matching orders
	matchingOrders, matchingErr := self.GetMatchingOrders(orderToInsert)

	if matchingErr != nil{
		return 0, matchingErr
	}

	//If there are no orders that match requested order, return nil
	if len(matchingOrders) == 0 {
		return insertId, nil
	}else{

		//Filter orders in the array
		orderToShareMap, filterError := self.FilterMatchingOrders(orderToInsert.NumShares, insertId, matchingOrders)

		if filterError != nil{
			fmt.Println(filterError.Error())
			if len(orderToShareMap) == 0{
				return 0, filterError
			}
		}

		//Create order transaction set
		orderTransactionSet := model.OrderTransactionPackage{
			MainOrder: insertId,
			MatchingOrders: orderToShareMap,
		}

		//Send order transaction package to concurrent transaction controller
		self.OrdersOutgoing <- orderTransactionSet
	}

	return insertId, nil
}


func (self *OrderController)GetMatchingOrders(order model.Order)([]model.Order, error){
	matchingOrders := make([]model.Order,0)

	sqlStatement := ""
	amountOfOrder, _ := order.CostPerShare.Float64()

	//If main order is a buy order, look for sell orders with price less than or equal to buy order price
	if order.InvestorAction == 0{
		sqlStatement = `SELECT * FROM orders WHERE symbol = ? AND investorAction = ? AND numSharesRemaining >= ? AND costPerShare <= ? AND (status = ? OR status = ?) ORDER BY created`
	//If main order is a sell order, look for buy orders with price greater than or equal to sell order price
	}else{
		sqlStatement = `SELECT * FROM orders WHERE symbol = ? AND investorAction = ? AND numSharesRemaining >= ? AND costPerShare >= ? AND (status = ? OR status = ?) ORDER BY created`

	}

	rows, err := self.Database.Query(sqlStatement, order.Symbol, util.IntAbsVal(int(order.InvestorAction)-1), 0, amountOfOrder, model.InProgress, model.Untouched)

	if err != nil{
		return matchingOrders, err
	}

	defer rows.Close()

	for rows.Next(){
		var (
			Id				int64
			Investor 		gosql.NullInt64
			Security 		gosql.NullInt64
			Symbol 			gosql.NullString
			InvestorAction	gosql.NullInt64
			InvestorType	gosql.NullInt64
			OrderType		gosql.NullInt64
			NumShares 		gosql.NullInt64
			NumSharesRemaining	gosql.NullInt64
			CostPerShare	gosql.NullFloat64
			CostOfShares	gosql.NullFloat64
			SystemFee 		gosql.NullFloat64
			TotalCost 		gosql.NullFloat64
			Created 		gosql.NullInt64
			Updated 		gosql.NullInt64
			Fulfilled 		gosql.NullInt64
			Status 			gosql.NullInt64
			AllowTakers 	gosql.NullBool
			LimitPerShare 	gosql.NullFloat64
			TakerFee		gosql.NullFloat64
			StopPrice		gosql.NullFloat64
		)

		err := rows.Scan(&Id , &Investor , &Security , &Symbol , &InvestorAction , &InvestorType , &OrderType , &NumShares ,
			&NumSharesRemaining, &CostPerShare , &CostOfShares , &SystemFee , &TotalCost , &Created , &Updated , &Fulfilled , &Status ,
			&AllowTakers , &LimitPerShare , &TakerFee, &StopPrice)

		if err != nil{
			return []model.Order{}, err
		}

		matchOrder := model.Order{
			ID: Id,
			Investor: Investor.Int64,
			Security: Security.Int64,
			Symbol: Symbol.String,
			InvestorAction: model.InvestorAction(InvestorAction.Int64),
			InvestorType: model.InvestorType(InvestorType.Int64),
			OrderType: model.OrderType(OrderType.Int64),
			NumShares: int(NumShares.Int64),
			NumSharesRemaining: int(NumSharesRemaining.Int64),
			CostPerShare: model.NewMoneyObject(CostPerShare.Float64),
			CostOfShares: model.NewMoneyObject(CostOfShares.Float64),
			SystemFee: model.NewMoneyObject(SystemFee.Float64),
			TotalCost: model.NewMoneyObject(TotalCost.Float64),
			Created: Created.Int64,
			Updated: Updated.Int64,
			Fulfilled: Fulfilled.Int64,
			Status: model.CompletionStatus(Status.Int64),
			AllowTakers: AllowTakers.Bool,
			LimitPerShare: model.NewMoneyObject(LimitPerShare.Float64),
			TakerFee: model.NewMoneyObject(LimitPerShare.Float64),
			StopPrice: model.NewMoneyObject(StopPrice.Float64),
		}

		matchingOrders = append(matchingOrders, matchOrder)
	}

	return matchingOrders, nil

}

func (self *OrderController) FilterMatchingOrders(currOrderNumShares int, currOrderID int64, matchingOrders []model.Order)(map[int64]int, error){
	//All orders sorted by timestamp in FIFO style
	orderToShareMap := make(map[int64]int, 0)

	numShares := currOrderNumShares

	workingNumShares := 0
	sqlQuery := `UPDATE orders SET status = ?, updated = ? WHERE id = ?`

	tx, err := self.Database.Begin()
	if err != nil{
		return map[int64]int{}, err
	}

	for _, order := range(matchingOrders){
		//If we haven't exceeded total amount of allotted shares
		if workingNumShares < numShares{
			//Calculate how many shares are left to be allotted
			sharesRemaining := numShares - workingNumShares

			//If the current order examined has more shares than we need allocated
			if order.NumSharesRemaining >= sharesRemaining{
				//Assign those shares
				orderToShareMap[order.ID] = sharesRemaining
				//Mark order as In Progress
				_, err := tx.Exec(sqlQuery, 1, time.Now().Unix(),order.ID)
				if err != nil{
					tx.Rollback()
					return map[int64]int{}, err
				}

				workingNumShares = numShares
				//Exit
				break
			//If current order examined has less shares than we need allocated
			}else{
				//Grab all the shares from this order
				orderToShareMap[order.ID] = order.NumSharesRemaining
				//Mark order as In Progress
				_, err := tx.Exec(sqlQuery, 1, time.Now().Unix(),order.ID)
				if err != nil{
					tx.Rollback()
					return map[int64]int{}, err
				}
				//Increment workingNumShares
				workingNumShares += order.NumSharesRemaining
				//Examine next order in next iteration of loop
			}
		}
	}

	if workingNumShares == 0{
		return map[int64]int{}, errors.New("NoMatchingOrder")
	}

	_, err = tx.Exec(sqlQuery, 1, time.Now().Unix(), currOrderID)
	if err != nil{
		tx.Rollback()
		return map[int64]int{}, err
	}

	err = tx.Commit()
	if err != nil{
		tx.Rollback()
		return map[int64]int{}, err
	}

	if numShares != workingNumShares{
		return orderToShareMap, errors.New("OrderPartiallyMatched")
	}

	return orderToShareMap, nil
}

func (self *OrderController) ValidateOrderFields(params model.OrderCreateParams) (error){
	//Field Validation
	if params.OrderType < 0 || params.InvestorType > 2{
		return errors.New("InvalidOrderType")
	}
	if params.InvestorType != 0 && params.InvestorType != 1{
		return errors.New("InvalidInvestorType")
	}
	if params.InvestorAction != 0 && params.InvestorAction != 1{
		return errors.New("InvalidInvestorAction")
	}

	//Verify Customer Exists in DB
	_, err := self.Database.GetCustomerByID(params.UserID)
	if err != nil {
		if err == gosql.ErrNoRows{
			return errors.New("CustomerDoesntExist")
		} else {
			return err
		}
	}

	//Verify symbol is valid
	securityID, err := self.Database.CheckIfSecurityInTable(params.Symbol)

	if err != nil {
		if err == gosql.ErrNoRows{
			return errors.New("SecurityDoesntExist")
		} else {
			return err
		}
	}

	//If Buy Action
	//Verify user has enough money to make purchase
	if params.InvestorAction == 0{
		portfolio, err := self.Database.GetPortfolioByCustomerID(params.UserID)
		if err != nil{
			self.Logger.Debug(err.Error())
			return errors.New("NoPortfolioForCustomer")
		}

		costOfShares := float64(params.NumShares) * params.CostPerShare
		totalAmount, _ := model.NewMoneyObject(costOfShares).Float64()
		portfolioCashValue, _ := portfolio.CashValue.Float64()

		//if !exactTotal{
		//	return errors.New("CostOfSharesUnExact")
		//}
		//
		//if !exactCash{
		//	return errors.New("PortfolioValueUnExact")
		//}

		if totalAmount > portfolioCashValue {
			return errors.New("InsufficientFunds")
		}
	}else{
		//If Sell Action
		//Make sure that user has enough of shares to satisfy transaction
		userOwnedShare, err := self.Database.GetOwnedShareForUserForSecurity(params.UserID, securityID)

		if err != nil{
			return errors.New("UserOwnedShareNonExistent")
		}

		if userOwnedShare.NumShares < params.NumShares{

			return errors.New("NotEnoughShares")
		}
	}


	return nil
}

func (self *OrderController)CreateAssocUserForEntity(entity model.Entity, params model.IPOParams)(int64, error){
	_, doesCustomerExist := self.Database.CheckIfCustomerInTable(entity.Email)

	if doesCustomerExist != nil {
		if doesCustomerExist == gosql.ErrNoRows{
		} else{
			return 0, doesCustomerExist
		}
	}

	security, err := self.Database.GetSecurityByID(entity.Security)
	if err != nil{
		return 0, err
	}

	//First insert user
	newUser := model.Customer{
		FirstName	: security.Symbol,
		LastName	: "Broker",
		PassHash	: entity.PassHash,
		Email 		: html.EscapeString(entity.Email),
		Portfolio	: 0,
	}

	userId, insertErr := self.Database.InsertCustomerIntoTable(newUser)

	if (insertErr != nil){
		return 0, insertErr
	}

	//Next insert new user portfolio
	//To start out with, each user will get $100.00 (Set at value and Cash Value)
	//This is solely for testing purposes
	//In the future, stripe integration will let us draw money from credit cards
	newUserPortfolio := model.Portfolio{
		Customer: userId,
		//Value = (Stock + Cash + Withdrawables)
		Value: model.NewMoneyObject(float64(params.NumShares) * params.SharePrice),
		StockValue: model.NewMoneyObject(float64(params.NumShares) * params.SharePrice),
		CashValue: model.NewMoneyObject(0.0),
		WithdrawableFunds: model.NewMoneyObject(0.0),
	}

	portfolioId, err := self.Database.InsertPortfolioToTable(newUserPortfolio)

	if err != nil{
		return 0, err
	}

	//Now update user to link it to the just created Portfolio
	updateErr := self.Database.AttachPortfolioToCustomer(userId, portfolioId)
	if updateErr != nil{
		return 0, err
	}

	ownedShareForIPO := model.OwnedShare{
		UserID: userId,
		Security: entity.Security,
		NumShares: params.NumShares,
	}

	_, ownedShareUpdateErr := self.Database.InsertOwnedShareToTable(ownedShareForIPO)
	if ownedShareUpdateErr != nil{
		return 0, err
	}

	return userId, nil
}

func (self *OrderController)GetMostOrderedSecuritiesOverTimeframe(startTime time.Time, numSecurities int)([]model.MostTradedReport, error){
	unixTime := startTime.Unix()

	sqlQuery := "SELECT `security`, COUNT(`security`) AS `value_occurrence` FROM `orders` WHERE created > ? GROUP BY `security` ORDER BY `value_occurrence` DESC LIMIT ?;"

	results, err := self.Database.Query(sqlQuery, unixTime, numSecurities)
	if err != nil{
		return []model.MostTradedReport{}, err
	}

	securities := make([]model.MostTradedReport, 0)

	for results.Next(){
		var securityId int64
		var frequency int64

		err := results.Scan(&securityId, &frequency)
		if err != nil{
			return []model.MostTradedReport{}, err
		}

		security, err := self.Database.GetSecurityByID(securityId)
		if err != nil{
			return []model.MostTradedReport{}, err
		}

		newMostTradedReport := model.MostTradedReport{
			Security: *security,
			Frequency: frequency,
		}

		securities = append(securities, newMostTradedReport)
	}

	return securities, nil
}


func (self *OrderController)IPO(params model.IPOParams, entityID int64)(int64, error){
	entity, err := self.Database.GetEntityByID(entityID)
	if err != nil{
		return 0, err
	}

	if entity.IPO == 1{
		return 0, errors.New(fmt.Sprintf("IPO for entity %d already happened", entityID))
	}

	if entity.Security == -1{
		return 0, errors.New(fmt.Sprintf("Entity %d does not have associated security", entityID))
	}

	security, err := self.Database.GetSecurityByEntityID(entityID)
	if err != nil{
		return 0, err
	}

	assocUser, assocUserErr := self.CreateAssocUserForEntity(*entity, params)

	if assocUserErr != nil{
		return 0, assocUserErr
	}

	updateErr := self.Database.UpdateEntityWithAssocUser(entityID, assocUser)
	if updateErr != nil{
		return 0, updateErr
	}


	IPOOrder := model.Order{
		Investor: assocUser,
		Security: entity.Security,
		Symbol: security.Symbol,
		InvestorAction: 1,
		InvestorType: 0,
		OrderType: 0,
		NumShares: params.NumShares,
		NumSharesRemaining: params.NumShares,
		CostPerShare: model.NewMoneyObject(params.SharePrice),
		CostOfShares: model.NewMoneyObject(params.SharePrice*float64(params.NumShares)),
		SystemFee: model.NewMoneyObject(0.0),
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
		Fulfilled: time.Now().Unix(),
		Status: model.CompletionStatus(0),
	}

	IPOOrder.TotalCost = model.NewMoneyObjectFromDecimal(IPOOrder.CostOfShares.Add(IPOOrder.SystemFee.Decimal))

	orderID, err := self.Database.InsertOrderIntoTable(IPOOrder)
	if err != nil{
		return 0, err
	}

	IPOTransaction := model.Transaction{
		OrderPlaced: orderID,
		NumShares: int64(params.NumShares),
		CostPerShare: IPOOrder.CostPerShare,
		TotalCost: model.NewMoneyObject(params.SharePrice*float64(params.NumShares)),
		SystemFee: model.NewMoneyObject(0.0),
		Created: time.Now().Unix(),
		Security: entity.Security,
	}

	ipoTransactionID, err := self.Database.InsertTransactionToTable(IPOTransaction)

	err = self.Database.CompleteEntityIPO(ipoTransactionID, entityID)
	if err != nil{
		return 0, err
	}

	return orderID, nil
}
//
//func (self *OrderController) InsertOrderToDatabase(params OrderCreateParams)(Order, error){
//	//Grab information about security
//	var security security.Security
//
//	securityCollection := mongo.GetSecurityCollection(mongo.GetDataBase(self.Session))
//
//	err := securityCollection.Find(bson.M{"symbol": params.Symbol}).One(&security)
//	if err != nil{
//		return Order{}, err
//	}
//
//	//Insert order
//	orderCollection := mongo.GetOrderBookForSecurity(mongo.GetDataBase(self.Session), params.Symbol)
//	costPerShare := api.NewMoneyObject(params.CostPerShare)
//	costOfShares := api.NewMoneyObject(params.SharesPurchased * params.CostPerShare)
//	systemFee := api.NewMoneyObject(0.0)
//	totalCost := costOfShares.Add(systemFee)
//
//	newOrder := Order{
//		Id: bson.NewObjectId(),
//		Investor: params.UserID,
//		Security: security.Id,
//		Symbol: security.Symbol,
//		Action: params.InvestorAction,
//		InvestorType: params.InvestorType,
//		OrderType: params.OrderType,
//		NumShares: params.SharesPurchased,
//		CostPerShare: costPerShare,
//		CostOfShares: costOfShares,
//		SystemFee: systemFee,
//		TotalCost: totalCost,
//		Created: time.Unix(0, params.TimeCreated),
//		Fulfilled: time.Unix(0, params.TimeCreated),
//		Updated: time.Unix(0, params.TimeCreated),
//		Status: api.Untouched,
//		Transactions: make([]bson.ObjectId, 0),
//		AllowTakers: false,
//		LimitPerShare: api.NewMoneyObject(0.0),
//		TakerFee: api.NewMoneyObject(0.0),
//		StopPrice: api.NewMoneyObject(0.0),
//	}
//
//	insertErr := orderCollection.Insert(newOrder)
//	if insertErr != nil{
//		return Order{}, insertErr
//	}
//	return newOrder, nil
//}
//
