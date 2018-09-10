package sql

import (
	"financialExchange/model"
	"encoding/json"
)

func (db *MySqlDB) CreateTransactionTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS transactions(" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"orderPlaced integer," +
		//ordersFulfilled is going to be a json array string of transaction IDs.
		//Works out bc transactions are immutable?
		"ordersFulfilled varChar(1000)," +
		"systemFee float," +
		"totalCost float," +
		"created int," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB) InsertTransactionToTable(transaction model.Transaction)(int64, error){
	query := `INSERT INTO transactions ( orderPlaced, ordersFulfilled, systemFee, totalCost, created ) VALUES ( ?, ?, ?, ?, ?)`

	ordersFulfilledJSONArray, _ := json.Marshal(transaction.OrdersFulfilled)

	r, err := db.Exec(query, transaction.OrderPlaced, string(ordersFulfilledJSONArray), transaction.SystemFee, transaction.TotalCost, transaction.Created)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()

	if err != nil{
		return 0, err
	}

	return lastInsertId, nil
}