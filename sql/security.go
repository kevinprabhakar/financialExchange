package sql

import (
	"financialExchange/model"
	"fmt"
	"database/sql"
)

func (db *MySqlDB) CreateSecurityTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS securities (" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"entity INT NULL," +
		"symbol varChar(10) NULL," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB)InsertSecurityIntoTable(security model.Security)(int64, error){
	query := `INSERT INTO securities (
		entity, symbol
	) VALUES (?, ?)`

	r, err := db.Exec(query, security.Entity, security.Symbol)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()
	if err != nil{
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *MySqlDB)CheckIfSecurityInTable(symbol string)(int64, error){
	query := `SELECT id FROM securities WHERE symbol = ?`

	row := db.QueryRow(query, symbol)

	var id int64

	err := row.Scan(&id)
	if err == sql.ErrNoRows{
		return 0, err
	}else{
		return id, err
	}

}