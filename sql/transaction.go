package sql

import (
	"financialExchange/model"
	gosql "database/sql"
)

func (db *MySqlDB) CreateTransactionTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS transactions(" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"orderPlaced integer," +
		//ordersFulfilled is going to be a json array string of transaction IDs.
		//Works out bc transactions are immutable?
		"numShares integer," +
		"costPerShare float," +
		"systemFee float," +
		"totalCost float," +
		"security integer," +
		"created int," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB) InsertTransactionToTable(transaction model.Transaction)(int64, error){
	query := `INSERT INTO transactions ( orderPlaced, numShares, costPerShare, systemFee, totalCost, security, created ) VALUES ( ?, ?, ?, ?, ?, ?, ?)`

	err := db.InsertFulfillingOrdersToTable(transaction.OrderPlaced, transaction.OrdersFulfilling)
	if err != nil{
		return 0, err
	}


	r, err := db.Exec(query, transaction.OrderPlaced, transaction.NumShares, transaction.CostPerShare, transaction.SystemFee, transaction.TotalCost, transaction.Security, transaction.Created)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()

	if err != nil{
		return 0, err
	}

	return lastInsertId, nil
}

func (db *MySqlDB) GetTransactionByID(transactionID int64)(model.Transaction, error){
	query := `SELECT * FROM transactions WHERE id = ?`

	rows, err := db.Query(query, transactionID)


	if err != nil{
		return model.Transaction{}, err
	}

	defer rows.Close()

	for rows.Next(){
		var (
			Id				int64
			OrderPlaced 		gosql.NullInt64
			NumShares 			gosql.NullInt64
			CostPerShare		gosql.NullFloat64
			SystemFee			gosql.NullFloat64
			TotalCost 			gosql.NullFloat64
			Security 			gosql.NullInt64
			Created 			gosql.NullInt64


		)

		err := rows.Scan(&Id , &OrderPlaced, &NumShares, &CostPerShare, &SystemFee, &TotalCost, &Security, &Created)

		if err != nil{
			return model.Transaction{}, err
		}

		newTransaction := model.Transaction{
			Id: Id,
			OrderPlaced: OrderPlaced.Int64,
			NumShares: NumShares.Int64,
			CostPerShare: model.NewMoneyObject(CostPerShare.Float64),
			SystemFee: model.NewMoneyObject(SystemFee.Float64),
			TotalCost: model.NewMoneyObject(TotalCost.Float64),
			Security: Security.Int64,
			Created: Created.Int64,
		}

		return newTransaction, nil

	}

	return model.Transaction{}, nil

}

func (db *MySqlDB) GetAllTransactionsForTimePeriodForSecurity(start int64, end int64, security int64)([]model.Transaction, error){
	query := `SELECT * FROM transactions WHERE created < ? AND created >= ? AND security = ?`

	rows, err := db.Query(query, end, start, security)


	if err != nil{
		return []model.Transaction{}, err
	}

	defer rows.Close()

	Transactions := make([]model.Transaction, 0)

	for rows.Next(){
		var (
			Id				int64
			OrderPlaced 		gosql.NullInt64
			NumShares 			gosql.NullInt64
			CostPerShare		gosql.NullFloat64
			SystemFee			gosql.NullFloat64
			TotalCost 			gosql.NullFloat64
			Security 			gosql.NullInt64
			Created 			gosql.NullInt64


		)

		err := rows.Scan(&Id , &OrderPlaced, &NumShares, &CostPerShare, &SystemFee, &TotalCost, &Security, &Created)

		if err != nil{
			return []model.Transaction{}, err
		}

		newTransaction := model.Transaction{
			Id: Id,
			OrderPlaced: OrderPlaced.Int64,
			NumShares: NumShares.Int64,
			CostPerShare: model.NewMoneyObject(CostPerShare.Float64),
			SystemFee: model.NewMoneyObject(SystemFee.Float64),
			TotalCost: model.NewMoneyObject(TotalCost.Float64),
			Security: Security.Int64,
			Created: Created.Int64,
		}

		Transactions = append(Transactions, newTransaction)
	}

	return Transactions, nil

}