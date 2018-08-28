package order

import (
	"gopkg.in/mgo.v2"
	"onePercent/util"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"financialExchange/mongo"
	"financialExchange/customer"
	"financialExchange/entity"
	"financialExchange/api"
	"financialExchange/portfolio"
	"math"
	"financialExchange/security"
	"time"
)

type OrderController struct {
	Session        *mgo.Session
	Logger         util.Logger
	OrdersOutgoing chan OrderTransactionPackage
}

func NewOrderController (session *mgo.Session, logger util.Logger, ordersOutgoing chan OrderTransactionPackage)(*OrderController){
	return &OrderController{
		Session: session,
		Logger: logger,
		OrdersOutgoing: ordersOutgoing,
	}
}

func (self *OrderController) CreateOrder(params OrderCreateParams)(error){
	//Validate order fields
	validErr := self.ValidateOrderFields(params)

	if validErr != nil{
		return validErr
	}

	//Insert order into database
	newOrder, insertErr := self.InsertOrderToDatabase(params)

	if insertErr != nil{
		return insertErr
	}

	//Get matching orders
	matchingOrders, matchingErr := self.GetMatchingOrders(newOrder, newOrder.Symbol)

	if matchingErr != nil{
		return matchingErr
	}

	//If there are no orders that match requested order, return nil
	if len(matchingOrders) == 0 {
		return nil
	}else{
		//Filter orders in the array
		orderToShareMap, filterError := self.FilterMatchingOrders(newOrder, matchingOrders)

		if filterError != nil{
			return filterError
		}

		//Create order transaction set
		orderTransactionSet := OrderTransactionPackage{
			MainOrder: newOrder.Id,
			MatchingOrders: orderToShareMap,
		}

		//Send order transaction package to concurrent transaction controller
		self.OrdersOutgoing <- orderTransactionSet

		return
	}


	return nil
}

func (self *OrderController)GetMatchingOrders(order Order, symbol string)([]Order, error){
	orderCollection := mongo.GetOrderBookForSecurity(mongo.GetDataBase(self.Session), symbol)

	matchingOrders := make([]Order, 0)


	//Only considering buy/sell signals..for now
	query := bson.M{
		"investorAction" : int(math.Abs(order.Action-1)),
		"numShares" : bson.M{
			"$gte" : order.NumShares,
		},
		"orderStatus" : bson.M{
			"$in" : []int{api.Untouched, api.InProgress},
		},
	}

	err := orderCollection.Find(query).Sort("-created").All(&matchingOrders)

	if err != nil{
		return []Order{}, err
	}

	return matchingOrders, nil

}

func (self *OrderController) FilterMatchingOrders(currOrder Order, matchingOrders []Order)(map[bson.ObjectId]int, error){
	//All orders sorted by timestamp in FIFO style
	orderToShareMap := make(map[bson.ObjectId]int, 0)

	numShares := currOrder.NumShares

	workingNumShares := 0


	for _, order := range(matchingOrders){
		//If we haven't exceeded total amount of allotted shares
		if workingNumShares < numShares{
			//Calculate how many shares are left to be allotted
			sharesRemaining := numShares - workingNumShares

			//If the current order examined has more shares than we need allocated
			if order.NumShares > sharesRemaining{
				//Assign those shares and exit
				orderToShareMap[order.Id] = sharesRemaining
				workingNumShares = numShares
				break
			//If current order examined has less shares than we need allocated
			}else{
				//Grab all the shares from this order
				orderToShareMap[order.Id] = order.NumShares
				//Increment workingNumShares
				workingNumShares += order.NumShares
				//Examine next order in next iteration of loop
			}
		}
	}

	if workingNumShares == 0{
		return map[bson.ObjectId]int{}, errors.New("NoMatchingOrder")
	}

	if numShares != workingNumShares{
		return orderToShareMap, errors.New("OrderPartiallyMatched")
	}

	return orderToShareMap, nil
}

func (self *OrderController) ValidateOrderFields(params OrderCreateParams) (error){
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

	//Verify User has signed up with this email before
	userCollection := mongo.GetCustomerCollection(mongo.GetDataBase(self.Session))
	portfolioCollection := mongo.GetPortfolioCollection(mongo.GetDataBase(self.Session))

	var findUser customer.Customer
	var userPortfolio portfolio.Portfolio

	err := userCollection.Find(bson.M{ "_id": bson.ObjectIdHex(params.UserID)}).One(&findUser)

	if (err == mgo.ErrNotFound){
		return errors.New("NonexistentUser")
	}
	if err != nil{
		return err
	}
	//Grab user portfolio
	err = portfolioCollection.Find(bson.M{"user":findUser.Id}).One(&userPortfolio)

	if (err == mgo.ErrNotFound){
		return errors.New("NonexistentUser")
	}
	if err != nil{
		return err
	}

	//Verify symbol is valid
	entityCollection := mongo.GetEntityCollection(mongo.GetDataBase(self.Session))

	var findEntity entity.Entity

	err = entityCollection.Find(bson.M{"symbol" : params.Symbol}).One(&findEntity)

	if (err == mgo.ErrNotFound){
		return errors.New("Nonexistent Symbol")
	}
	if err != nil{
		return err
	}

	//Verify user has enough money to make purchase
	costOfShares := params.SharesPurchased * params.CostPerShare
	totalAmount := api.NewMoneyObject(costOfShares)

	if totalAmount > userPortfolio.CashValue {
		return errors.New("InsufficientFunds")
	}

	return nil
}

func (self *OrderController) InsertOrderToDatabase(params OrderCreateParams)(Order, error){
	//Grab information about security
	var security security.Security

	securityCollection := mongo.GetSecurityCollection(mongo.GetDataBase(self.Session))

	err := securityCollection.Find(bson.M{"symbol": params.Symbol}).One(&security)
	if err != nil{
		return Order{}, err
	}

	//Insert order
	orderCollection := mongo.GetOrderBookForSecurity(mongo.GetDataBase(self.Session), params.Symbol)
	costPerShare := api.NewMoneyObject(params.CostPerShare)
	costOfShares := api.NewMoneyObject(params.SharesPurchased * params.CostPerShare)
	systemFee := api.NewMoneyObject(0.0)
	totalCost := costOfShares.Add(systemFee)

	newOrder := Order{
		Id: bson.NewObjectId(),
		Investor: params.UserID,
		Security: security.Id,
		Symbol: security.Symbol,
		Action: params.InvestorAction,
		InvestorType: params.InvestorType,
		OrderType: params.OrderType,
		NumShares: params.SharesPurchased,
		CostPerShare: costPerShare,
		CostOfShares: costOfShares,
		SystemFee: systemFee,
		TotalCost: totalCost,
		Created: time.Unix(0, params.TimeCreated),
		Fulfilled: time.Unix(0, params.TimeCreated),
		Updated: time.Unix(0, params.TimeCreated),
		Status: api.Untouched,
		Transactions: make([]bson.ObjectId, 0),
		AllowTakers: false,
		LimitPerShare: api.NewMoneyObject(0.0),
		TakerFee: api.NewMoneyObject(0.0),
		StopPrice: api.NewMoneyObject(0.0),
	}

	insertErr := orderCollection.Insert(newOrder)
	if insertErr != nil{
		return Order{}, insertErr
	}
	return newOrder, nil
}

