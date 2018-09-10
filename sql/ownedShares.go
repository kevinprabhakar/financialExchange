package sql

import (
	"financialExchange/model"
	gosql "database/sql"
	"fmt"
)

func (db *MySqlDB) CreateOwnedSharesTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS ownedShares (" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"userId INT NULL," +
		"security INT NULL," +
		"numShares INT NULL," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB) InsertOwnedShareToTable(ownedShare model.OwnedShare)(int64, error){
	query := `INSERT INTO ownedShares (userId, security, numShares) VALUES ( ?, ?, ? )`

	r, err := db.Exec(query, ownedShare.UserID, ownedShare.Security, ownedShare.NumShares)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()

	if err != nil{
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func ScanOwnedShare(s RowScanner)(model.OwnedShare, error){
	var(
		Id 			gosql.NullInt64
		UserID		gosql.NullInt64
		Security	gosql.NullInt64
		NumShares 	gosql.NullInt64
	)

	if err := s.Scan(&Id, &UserID, &Security, &NumShares); err != nil{
		return model.OwnedShare{}, err
	}

	ownedShare := model.OwnedShare{
		UserID: UserID.Int64,
		Security: Security.Int64,
		NumShares: int(NumShares.Int64),
	}

	return ownedShare, nil
}

func (db *MySqlDB)GetOwnedShareForUserForSecurity(userID int64, securityId int64)(model.OwnedShare, error){
	query := `SELECT * FROM ownedShares WHERE userId = ?, security = ?`

	row := db.QueryRow(query, userID, securityId)

	ownedShare, err := ScanOwnedShare(row)
	if err != nil{
		return model.OwnedShare{}, err
	}
	return ownedShare, nil
}

func (db *MySqlDB) GetAllOwnedSharesForUserID(userID int64)([]model.OwnedShare, error){
	query := `SELECT * FROM ownedShares WHERE userId = ?`

	rows, err := db.Query(query, userID)

	if err != nil{
		return []model.OwnedShare{}, err
	}

	defer rows.Close()

	OwnedShares := make([]model.OwnedShare, 0)

	for rows.Next(){
		var (
			Id				gosql.NullInt64
			Investor 		gosql.NullInt64
			Security 		gosql.NullInt64
			NumShares 		gosql.NullInt64
		)

		err := rows.Scan(&Id , &Investor , &Security , &NumShares)

		if err != nil{
			return []model.OwnedShare{}, err
		}

		matchOrder := model.OwnedShare{
			UserID: Investor.Int64,
			Security: Security.Int64,
			NumShares: int(NumShares.Int64),
		}

		OwnedShares = append(OwnedShares, matchOrder)
	}

	return OwnedShares, nil

}