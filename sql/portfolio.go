package sql

import (
	"financialExchange/model"
	"database/sql"
)

func (db *MySqlDB) CreatePortfolioTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS portfolios (" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"customer integer," +
		"value float," +
		"stockValue float," +
		"cashValue float," +
		"withdrawableFunds float," +
		"PRIMARY KEY (id)" +
		" )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB) InsertPortfolioToTable(portfolio model.Portfolio)(int64, error){
	query := `INSERT INTO portfolios (
		customer, value, stockValue, cashValue, withdrawableFunds
	) VALUES (?, ?, ?, ?, ?)`

	r, err := db.Exec(query, portfolio.Customer, portfolio.Value, portfolio.StockValue, portfolio.CashValue, portfolio.WithdrawableFunds)
	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()
	if err != nil{
		return 0, err
	}

	return lastInsertId, nil
}

func ScanPortfolio(s RowScanner)(*model.Portfolio, error){
	var (
		id 		int64
		Customer sql.NullInt64
		Value	sql.NullFloat64
		StockValue 	sql.NullFloat64
		CashValue sql.NullFloat64
		WithdrawableFunds sql.NullFloat64
	)
	if err := s.Scan(&id, &Customer, &Value, &StockValue, &CashValue, &WithdrawableFunds); err != nil{
		return nil, err
	}

	portfolio := &model.Portfolio{
		Customer: Customer.Int64,
		Value: model.NewMoneyObject(Value.Float64),
		StockValue: model.NewMoneyObject(StockValue.Float64),
		CashValue: model.NewMoneyObject(CashValue.Float64),
		WithdrawableFunds: model.NewMoneyObject(WithdrawableFunds.Float64),
	}

	return portfolio, nil
}

func (db *MySqlDB)GetPortfolioByCustomerID(customerID int64)(*model.Portfolio, error){
	sqlStatement := `SELECT * FROM portfolios WHERE customer = ?`

	row := db.QueryRow(sqlStatement, customerID)

	user, err := ScanPortfolio(row)
	if err != nil{
		return nil, err
	}
	return user, nil
}