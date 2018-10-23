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
		"created integer," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB)InsertSecurityIntoTable(security model.Security)(int64, error){
	query := `INSERT INTO securities (
		entity, symbol, created
	) VALUES (?, ?, ?)`

	r, err := db.Exec(query, security.Entity, security.Symbol, security.Created)

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

func ScanSecurity(s RowScanner)(*model.Security, error){
	var (
		id 		int64
		entity 	sql.NullInt64
		symbol sql.NullString
		created sql.NullInt64
	)
	if err := s.Scan(&id, &entity, &symbol, &created); err != nil{
		return nil, err
	}

	security := &model.Security{
		Id: id,
		Entity: entity.Int64,
		Symbol: symbol.String,
		Created: created.Int64,
	}

	return security, nil
}

func (db *MySqlDB)GetSecurityByID(uid int64)(*model.Security, error){
	sqlStatement := `SELECT * FROM securities WHERE id = ?`

	row := db.QueryRow(sqlStatement, uid)

	security, err := ScanSecurity(row)
	if err != nil{
		return nil, err
	}
	return security, nil
}

func (db *MySqlDB)GetSecurityByEntityID(uid int64)(*model.Security, error){
	sqlStatement := `SELECT * FROM securities WHERE entity = ?`

	row := db.QueryRow(sqlStatement, uid)

	security, err := ScanSecurity(row)
	if err != nil{
		return nil, err
	}
	return security, nil
}