package sql

import (
	"fmt"
	"financialExchange/model"
	"database/sql"
)

func (db *MySqlDB) CreateOrderTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS orders(" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"investor integer," +
		"security integer," +
		"symbol varChar(4)," +
		"investorAction integer," +
		"investorType integer," +
		"orderType integer," +
		"numShares integer," +
		"costPerShare float," +
		"costOfShares float," +
		"systemFee float," +
		"totalCost float," +
		"created int," +
		"updated int," +
		"fulfilled int," +
		"status	integer," +
		"allowTakers boolean," +
		"limitPerShare float," +
		"takerFee float," +
		"stopPrice float,"+
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB)InsertOrderIntoTable(order model.Order) (int64, error){
	query := `INSERT INTO orders (
		investor, security, symbol, investorAction, investorType, orderType, numShares, costPerShare, costOfShares, systemFee,
		totalCost, created, updated, fulfilled, status, allowTakers, limitPerShare, takerFee, stopPrice
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	r, err := db.Exec(query, order.Investor, order.Security, order.Symbol, order.InvestorAction, order.InvestorType, order.OrderType,
			order.NumShares, order.CostPerShare, order.CostOfShares, order.SystemFee, order.TotalCost, order.Created, order.Updated, order.Fulfilled,
			order.Status, order.AllowTakers, order.LimitPerShare, order.TakerFee, order.StopPrice)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()
	if err != nil{
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}


func ScanOrder(s RowScanner)(*model.Order, error){
	var (
		Id				sql.NullInt64
		Investor 		sql.NullInt64
		Security 		sql.NullInt64
		Symbol 			sql.NullString
		InvestorAction	sql.NullInt64
		InvestorType	sql.NullInt64
		OrderType		sql.NullInt64
		NumShares 		sql.NullInt64
		CostPerShare	sql.NullFloat64
		CostOfShares	sql.NullFloat64
		SystemFee 		sql.NullFloat64
		TotalCost 		sql.NullFloat64
		Created 		sql.NullInt64
		Updated 		sql.NullInt64
		Fulfilled 		sql.NullInt64
		Status 			sql.NullInt64
		AllowTakers 	sql.NullBool
		LimitPerShare 	sql.NullFloat64
		TakerFee		sql.NullFloat64
		StopPrice		sql.NullFloat64
	)
	if err := s.Scan(&Id , &Investor , &Security , &Symbol , &InvestorAction , &InvestorType , &OrderType , &NumShares ,
		&CostPerShare , &CostOfShares , &SystemFee , &TotalCost , &Created , &Updated , &Fulfilled , &Status ,
		&AllowTakers , &LimitPerShare , &TakerFee, &StopPrice); err != nil{
		return nil, err
	}

	order := &model.Order{
		ID: Id.Int64,
		Investor: Investor.Int64,
		Security: Security.Int64,
		Symbol: Symbol.String,
		InvestorAction: model.InvestorAction(InvestorAction.Int64),
		InvestorType: model.InvestorType(InvestorType.Int64),
		OrderType: model.OrderType(OrderType.Int64),
		NumShares: int(NumShares.Int64),
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

	return order, nil
}

func (db *MySqlDB)GetOrderById(id int64)(*model.Order, error){
	sqlStatement := `SELECT * FROM orders WHERE id = ?`

	row := db.QueryRow(sqlStatement, id)

	user, err := ScanOrder(row)
	if err != nil{
		return nil, err
	}
	return user, nil
}
