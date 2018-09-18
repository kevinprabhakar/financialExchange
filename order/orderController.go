package order

import (
	"financialExchange/util"
	"financialExchange/model"
	"financialExchange/sql"
	"errors"
	gosql "database/sql"
	"time"

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
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
		Fulfilled: time.Now().Unix(),
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
		orderToShareMap, filterError := self.FilterMatchingOrders(orderToInsert, matchingOrders)

		if filterError != nil{
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

	sqlStatement := `SELECT * FROM orders WHERE symbol = ? AND investorAction = ? AND numSharesRemaining >= ? AND (status = ? OR status = ?) ORDER BY created`
	rows, err := self.Database.Query(sqlStatement, order.Symbol, util.IntAbsVal(int(order.InvestorAction)-1), 0, model.InProgress, model.Untouched)

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

func (self *OrderController) FilterMatchingOrders(currOrder model.Order, matchingOrders []model.Order)(map[int64]int, error){
	//All orders sorted by timestamp in FIFO style
	orderToShareMap := make(map[int64]int, 0)

	numShares := currOrder.NumShares

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

	if numShares != workingNumShares{
		return orderToShareMap, errors.New("OrderPartiallyMatched")
	}

	_, err = tx.Exec(sqlQuery, 1, time.Now().Unix(), currOrder.ID)
	if err != nil{
		tx.Rollback()
		return map[int64]int{}, err
	}

	err = tx.Commit()
	if err != nil{
		tx.Rollback()
		return map[int64]int{}, err
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
		totalAmount, exactTotal := model.NewMoneyObject(costOfShares).Float64()
		portfolioCashValue, exactCash := portfolio.CashValue.Float64()

		if !exactTotal{
			return errors.New("CostOfSharesUnExact")
		}

		if !exactCash{
			return errors.New("PortfolioValueUnExact")
		}

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
