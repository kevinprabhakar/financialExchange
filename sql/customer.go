package sql

import (
	"fmt"
	"database/sql"
	"financialExchange/model"

)

func (db *MySqlDB) CreateCustomerTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS customers (" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"firstName varchar(100) NULL," +
		"lastName varChar(100) NULL," +
		"email varChar(100) NULL," +
		"passHash varChar(400) NULL," +
		"portfolio integer NULL," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB)InsertCustomerIntoTable(customer model.Customer)(int64, error){
	query := `INSERT INTO customers (
		firstName, lastName, email, passHash, portfolio
	) VALUES (?, ?, ?, ?, ?)`

	r, err := db.Exec(query, customer.FirstName, customer.LastName, customer.Email, customer.PassHash, customer.Portfolio)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()
	if err != nil{
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *MySqlDB)AttachPortfolioToCustomer(userId int64, portfolioID int64)(error) {
	sqlStatement := fmt.Sprint("UPDATE customers SET portfolio = ? WHERE id = ?")

	_, err:= db.Exec(sqlStatement, portfolioID, userId)

	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB)CheckIfCustomerInTable(email string)(string, error) {
	sqlStatement := fmt.Sprint("SELECT email FROM customers where email = ?")

	row := db.QueryRow(sqlStatement, email)

	var emailCheck string

	err := row.Scan(&emailCheck)
	if err == sql.ErrNoRows{
		return "", err
	}else{
		return emailCheck, err
	}
}

func ScanCustomer(s RowScanner)(*model.Customer, error){
	var (
		id 		int64
		firstName 	sql.NullString
		lastName sql.NullString
		passHash sql.NullString
		email    sql.NullString
		portfolio int64
	)
	if err := s.Scan(&id, &firstName, &lastName, &email, &passHash, &portfolio); err != nil{
		return nil, err
	}

	customer := &model.Customer{
		Id: id,
		FirstName: firstName.String,
		LastName: lastName.String,
		PassHash: passHash.String,
		Email: email.String,
		Portfolio: portfolio,
	}

	return customer, nil
}

func (db *MySqlDB)GetCustomerByEmail(email string)(*model.Customer, error){
	sqlStatement := `SELECT * FROM customers WHERE email = ?`

	row := db.QueryRow(sqlStatement, email)

	user, err := ScanCustomer(row)
	if err != nil{
		return nil, err
	}
	return user, nil
}

func (db *MySqlDB)GetCustomerByID(ID int64)(*model.Customer, error){
	sqlStatement := `SELECT * FROM customers WHERE id = ?`

	row := db.QueryRow(sqlStatement, ID)

	user, err := ScanCustomer(row)
	if err != nil{
		return nil, err
	}
	return user, nil
}