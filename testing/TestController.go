package testing

import (
	"financialExchange/sql"
	"financialExchange/util"
	"financialExchange/model"
	"time"
	"math/rand"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"errors"
	"strconv"
	"financialExchange/order"
	"math"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")


type TestController struct{
	Db 		*sql.MySqlDB
	Logger 	*util.Logger
}

func NewTestController(db *sql.MySqlDB, logger *util.Logger)(*TestController){
	return &TestController{
		Db: db,
		Logger: logger,
	}
}

func RandSeq(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (self *TestController)DropAllTables()(error){
	tableNames := []string{"customers","orders","entities","ownedShares","portfolios","securities","transactions","fulfillingOrders"}

	tx, err := self.Db.Begin()
	if err != nil{
		self.Logger.ErrorMsg("Error creating drop table transaction")
		return err
	}

	for _, tableName := range(tableNames){
		sqlQuery := fmt.Sprintf("DROP TABLE %s", tableName)
		_, err := tx.Exec(sqlQuery)

		if err != nil{
			self.Logger.ErrorMsg("Error Executing DROP TABLE transaction")
			tx.Rollback()
			self.Logger.Debug(err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil{
		self.Logger.ErrorMsg("Error committing drop table transaction")
		return err
	}

	return nil
}

func (self *TestController)CreateAllTables()(error){
	err := self.Db.CreateCustomerTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Customer Table")
		return err
	}

	err = self.Db.CreatePortfolioTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Portfolio Table")

		return err
	}

	err = self.Db.CreateEntityTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Entity Table")

		return err
	}

	err = self.Db.CreateSecurityTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Security Table")

		return err
	}

	err = self.Db.CreateOrderTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Order Table")

		return err
	}

	err = self.Db.CreateOwnedSharesTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Owned Shares Table")


		return err
	}

	err = self.Db.CreateTransactionTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating Transaction Table")

		return err
	}
	err = self.Db.CreateFulfillingOrdersTable()
	if err != nil{
		self.Logger.ErrorMsg("Error Creating FulfillingOrders Table")

		return err
	}
	return nil
}

func (self *TestController)ResetTables()(error){
	err := self.DropAllTables()
	if err != nil{
		return err
	}

	err = self.CreateAllTables()
	if err != nil{
		return err
	}
	return nil
}

//All functions from this point out assume clean tables
func (self *TestController)CreateUser()(int64, error){
	randPassword, _ := util.HashPassword(RandSeq(10))
	newCustomer := model.Customer{
		Email: RandSeq(10)+"@gmail.com",
		PassHash: randPassword,
		FirstName: RandSeq(10),
		LastName: RandSeq(10),
		Portfolio: 0,
	}

	userId, err := self.Db.InsertCustomerIntoTable(newCustomer)
	if err != nil{
		self.Logger.ErrorMsg("Error inserting customer into table")
		return 0, err
	}

	newUserPortfolio := model.Portfolio{
		Customer: userId,
		Value: model.NewMoneyObject(0.0),
		StockValue: model.NewMoneyObject(0.0),
		CashValue: model.NewMoneyObject(0.0),
		WithdrawableFunds: model.NewMoneyObject(0.0),
	}

	portfolioId, err := self.Db.InsertPortfolioToTable(newUserPortfolio)
	if err != nil{
		self.Logger.ErrorMsg("Error inserting portfolio into table")
		return 0, err
	}

	err = self.Db.AttachPortfolioToCustomer(userId, portfolioId)
	if err != nil{
		self.Logger.ErrorMsg("Error attaching portfolio to customer")

		return 0, err
	}
	return userId, nil
}

func (self *TestController)GiveUserMoney(userId int64, cashValue model.Money)(error){
	sqlQuery := `UPDATE portfolios SET value = ?, cashValue = ? WHERE id = ?`

	currPortfolio, err := self.Db.GetPortfolioByCustomerID(userId)
	if err != nil{
		self.Logger.ErrorMsg("Error getting portfolio for customer")

		return err
	}

	currCashValue, _ := currPortfolio.CashValue.Float64()
	currValue, _ := currPortfolio.Value.Float64()
	newCashValue, _ := cashValue.Float64()

	_, err = self.Db.Exec(sqlQuery, currValue+newCashValue, currCashValue+newCashValue, userId)

	if err != nil{
		self.Logger.ErrorMsg("Error giving customer money")
		self.Logger.Debug(err.Error())
		return err
	}

	return nil
}

func (self *TestController)CreateEntityWithSecurityWithNameAndSymbol(entityName string, entitySymbol string)(int64, int64, string, error){
	newEntity := model.Entity{
		Name: entityName,
		Email: RandSeq(10)+"@gmail.com",
		PassHash: RandSeq(10),
		Security: -1,
		Created: time.Now().Unix(),
		Deleted: time.Now().Unix(),
	}

	entityId, insertErr := self.Db.InsertEntityIntoTable(newEntity)

	if insertErr != nil{
		self.Logger.ErrorMsg("Error inserting entity into table")
		return 0, 0, "", insertErr
	}

	//Create security and insert to database
	newSecurity := model.Security{
		Entity: entityId,
		Symbol: entitySymbol,
	}

	securityId, insertErr := self.Db.InsertSecurityIntoTable(newSecurity)

	if insertErr != nil {
		self.Logger.ErrorMsg("Error inserting security into table")
		return 0, 0, "", insertErr
	}

	err := self.Db.UpdateEntityWithSecurityID(entityId, securityId)
	if err != nil{
		self.Logger.ErrorMsg("Error updating entity with security Id")
		return 0, 0, "", err
	}

	return entityId, securityId, newSecurity.Symbol, nil
}

//Returns entity ID, security Id, security Symbol, and error
func (self *TestController)CreateEntityWithSecurity()(int64, int64, string, error){
	newEntity := model.Entity{
		Name: RandSeq(10),
		Email: RandSeq(10)+"@gmail.com",
		PassHash: RandSeq(10),
		Security: -1,
		Created: time.Now().Unix(),
		Deleted: time.Now().Unix(),
	}

	entityId, insertErr := self.Db.InsertEntityIntoTable(newEntity)

	if insertErr != nil{
		self.Logger.ErrorMsg("Error inserting entity into table")
		return 0, 0, "", insertErr
	}

	//Create security and insert to database
	newSecurity := model.Security{
		Entity: entityId,
		Symbol: RandSeq(3),
	}

	securityId, insertErr := self.Db.InsertSecurityIntoTable(newSecurity)

	if insertErr != nil {
		self.Logger.ErrorMsg("Error inserting security into table")
		return 0, 0, "", insertErr
	}

	err := self.Db.UpdateEntityWithSecurityID(entityId, securityId)
	if err != nil{
		self.Logger.ErrorMsg("Error updating entity with security Id")
		return 0, 0, "", err
	}

	return entityId, securityId, newSecurity.Symbol, nil
}

func (self *TestController)IPOEntity(entityID int64)(int64, error){
	newChannel := make(chan model.OrderTransactionPackage, 0)
	OrderController := order.NewOrderController(self.Db, self.Logger, newChannel)
	IPOParams := model.IPOParams{
		SharePrice: 100.00,
		NumShares: 40000,
	}

	orderID, err := OrderController.IPO(IPOParams, entityID)
	if err != nil{
		return 0, err
	}

	return orderID, nil
}

func (self *TestController)GiveUserStock(userId int64, security int64, numShares int, price model.Money)(error){
	newOwnedShare := model.OwnedShare{
		UserID: userId,
		Security: security,
		NumShares: numShares,
	}
	_, err := self.Db.InsertOwnedShareToTable(newOwnedShare)
	if err != nil{
		self.Logger.ErrorMsg("Error inserting owned share to table")
		return err
	}

	currPortfolio, err := self.Db.GetPortfolioByCustomerID(userId)
	if err != nil{
		self.Logger.ErrorMsg("Error getting customer portfolio")
		return nil
	}
	priceFloat, _ := price.Float64()
	currValue, _ := currPortfolio.Value.Float64()
	currStockValue, _ := currPortfolio.StockValue.Float64()

	sqlQuery := `UPDATE portfolios SET value = ?, stockValue = ? WHERE id = ?`
	_ , err = self.Db.Exec(sqlQuery, currValue+(priceFloat*float64(numShares)), currStockValue+(priceFloat*float64(numShares)), userId)
	if err != nil{
		self.Logger.ErrorMsg("Error updating portfolio value for customer")
		self.Logger.Debug(err.Error())
		return err
	}
	return nil
}


func (self *TestController)PlaceBuyOrderForUser(userId int64, security int64, symbol string, price model.Money, numShares int)(int64, error){
	stockPrice, _ := price.Float64()

	url := "http://localhost:3001/api/order"

	newOrderParams := model.OrderCreateParams{
		UserID: userId,
		InvestorAction: 0,
		InvestorType: 0,
		OrderType: 0,
		Symbol: symbol,
		NumShares: numShares,
		CostPerShare: stockPrice,
		TimeCreated: time.Now().Unix(),
		AllowTakers: true,
		LimitPerShare: 0.0,
		StopPrice: 0.0,
	}

	byteForm, err := json.Marshal(newOrderParams)
	if err != nil{
		self.Logger.ErrorMsg("Error decoding JSON order params")
		return 0, err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteForm))

	accessToken, err := util.GetAccessToken(strconv.FormatInt(userId, 10))
	if err != nil{
		self.Logger.ErrorMsg("Error getting access token for user")
		return 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil{
		self.Logger.ErrorMsg("Error placing buy order")
		return 0, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200{
		return 0, errors.New("Error: "+ string(body))
	}

	var orderIdStruct model.OrderId

	err = json.Unmarshal(body, &orderIdStruct)
	if err != nil{
		return 0, err
	}

	return orderIdStruct.Id, nil
}

func (self *TestController)PlaceBuyOrderForUserWithSpecTime(userId int64, security int64, symbol string, price model.Money, numShares int, timeCreated int64)(int64, error){
	stockPrice, _ := price.Float64()

	url := "http://localhost:3001/api/order"

	newOrderParams := model.OrderCreateParams{
		UserID: userId,
		InvestorAction: 0,
		InvestorType: 0,
		OrderType: 0,
		Symbol: symbol,
		NumShares: numShares,
		CostPerShare: stockPrice,
		TimeCreated: timeCreated,
		AllowTakers: true,
		LimitPerShare: 0.0,
		StopPrice: 0.0,
	}

	byteForm, err := json.Marshal(newOrderParams)
	if err != nil{
		self.Logger.ErrorMsg("Error decoding JSON order params")
		return 0, err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteForm))

	accessToken, err := util.GetAccessToken(strconv.FormatInt(userId, 10))
	if err != nil{
		self.Logger.ErrorMsg("Error getting access token for user")
		return 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil{
		self.Logger.ErrorMsg("Error placing buy order")
		return 0, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200{
		return 0, errors.New("Error: "+ string(body))
	}

	var orderIdStruct model.OrderId

	err = json.Unmarshal(body, &orderIdStruct)
	if err != nil{
		return 0, err
	}

	return orderIdStruct.Id, nil
}

func (self *TestController)PlaceSellOrderForUser(userId int64, security int64, symbol string, price model.Money, numShares int)(int64, error){
	stockPrice, _ := price.Float64()

	url := "http://localhost:3001/api/order"

	newOrderParams := model.OrderCreateParams{
		UserID: userId,
		InvestorAction: 1,
		InvestorType: 0,
		OrderType: 0,
		Symbol: symbol,
		NumShares: numShares,
		CostPerShare: stockPrice,
		TimeCreated: time.Now().Unix(),
		AllowTakers: true,
		LimitPerShare: 0.0,
		StopPrice: 0.0,
	}

	byteForm, err := json.Marshal(newOrderParams)
	if err != nil{
		self.Logger.ErrorMsg("Error decoding JSON order params")
		return 0, err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteForm))

	accessToken, err := util.GetAccessToken(strconv.FormatInt(userId, 10))
	if err != nil{
		self.Logger.ErrorMsg("Error getting access token for user")
		return 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil{
		self.Logger.ErrorMsg("Error placing sell order")
		return 0, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200{
		return 0, errors.New("Error: "+ string(body))
	}

	var orderIdStruct model.OrderId

	err = json.Unmarshal(body, &orderIdStruct)
	if err != nil{
		return 0, err
	}

	return orderIdStruct.Id, nil
}

func(self *TestController)Brownian(seed int64, N float64)([]float64, []float64){
	r2 := rand.New(rand.NewSource(seed))
	dt := 1./N
	bVals := make([]float64, 0)
	totalSum := make([]float64, 0)
	for i := 0; i < int(N); i++{
		b := r2.NormFloat64() * math.Sqrt(dt)
		if len(totalSum) == 0{
			totalSum = append(totalSum, b)
		}else{
			totalSum = append(totalSum, totalSum[len(totalSum)-1]+b)
		}
		bVals = append(bVals, b)
	}
	bVals = append([]float64{0.0}, bVals...)
	return totalSum, bVals
}

func(self *TestController)GeometricBrownianMotion(s0 float64, mu float64, sigma float64, t float64, N float64)([]float64, []float64, error){
	linSpace := make([]float64, 0)
	for i := float64(0); i <= N; i += 1 {
		linSpace = append(linSpace, i/N)
	}
	S := make([]float64, 0)
	S = append(S, s0)

	W, _ := self.Brownian(time.Now().Unix(), N)

	for i := float64(1); i <= N; i++{
		drift := (mu - 0.5 * math.Pow(sigma, 2.0)) * linSpace[int(i)]
		diffusion := sigma * W[int(i)-1]
		S_temp := s0 * math.Exp(drift + diffusion)
		S = append(S, S_temp)
	}

	return S, linSpace, nil
}

func(self *TestController)CreateTransactionsForEntityFromGBM(entityID int64, prices []float64)(error){
	startTime := time.Now().Add(-1 * 24* time.Hour)
	timeIncrement := int64(time.Now().Sub(startTime).Seconds()/float64(len(prices)))
	entity, err := self.Db.GetEntityByID(entityID)
	if err != nil{
		return err
	}


	for index, price := range prices{
		tempTransaction := model.Transaction{
			OrderPlaced: int64(index),
			OrdersFulfilling: []int64{0},
			NumShares: int64(1),
			CostPerShare: model.NewMoneyObject(price),
			TotalCost: model.NewMoneyObject(price),
			SystemFee: model.NewMoneyObject(0.0),
			Created: int64((int64(index)+1)*timeIncrement)+startTime.Unix(),
			Security: entity.Security,

		}
		_, err = self.Db.InsertTransactionToTable(tempTransaction)
		if err != nil{
			return err
		}
	}
	return nil
}