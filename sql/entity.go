package sql

import (
	"financialExchange/model"
	"fmt"
	"database/sql"
)

func (db *MySqlDB) CreateEntityTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS entities (" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"name varchar(100) NULL," +
		"email varChar(100) NULL," +
		"passHash varChar(400) NULL," +
		"security BIGINT NULL," +
		"created BIGINT NULL," +
		"deleted BIGINT NULL," +
		"ipo BIGINT NULL," +
		"assocUser INT," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB)InsertEntityIntoTable(entity model.Entity)(int64, error){
	query := `INSERT INTO entities (
		name, email, passHash, security, created, deleted, ipo, assocUser
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	r, err := db.Exec(query, entity.Name, entity.Email, entity.PassHash, entity.Security, entity.Created, entity.Deleted, entity.IPO, entity.AssocUser)

	if err != nil{
		return 0, err
	}

	lastInsertId, err := r.LastInsertId()

	if err != nil{
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertId, nil
}

func (db *MySqlDB)CheckIfEntityInTable(email string)(string, error) {
	sqlStatement := fmt.Sprint("SELECT * FROM entities where email = ?")

	row := db.QueryRow(sqlStatement, email)

	var emailCheck string

	err := row.Scan(&emailCheck)
	if err == sql.ErrNoRows{
		return "", err
	}else{
		return emailCheck, err
	}
}

func ScanEntity(s RowScanner)(*model.Entity, error){
	var (
		id 		int64
		name 	sql.NullString
		email sql.NullString
		passHash sql.NullString
		security    int64
		created int64
		deleted int64
		ipo 	int64
		assocUser int64
	)
	if err := s.Scan(&id, &name, &email, &passHash, &security, &created, &deleted, &ipo, &assocUser); err != nil{
		return nil, err
	}

	entity := &model.Entity{
		Id: id,
		Name: name.String,
		Email: email.String,
		PassHash: passHash.String,
		Security: security,
		Created: created,
		Deleted: deleted,
		IPO: ipo,
		AssocUser: assocUser,
	}

	return entity, nil
}

func (db *MySqlDB)GetEntityByEmail(email string)(*model.Entity, error){
	sqlStatement := `SELECT * FROM entities WHERE email = ?`

	row := db.QueryRow(sqlStatement, email)

	user, err := ScanEntity(row)
	if err != nil{
		return nil, err
	}
	return user, nil
}

func (db *MySqlDB)GetEntityByID(uid int64)(*model.Entity, error){
	sqlStatement := `SELECT * FROM entities WHERE id = ?`

	row := db.QueryRow(sqlStatement, uid)

	user, err := ScanEntity(row)
	if err != nil{
		return nil, err
	}
	return user, nil
}

func (db *MySqlDB)UpdateEntityWithSecurityID(entityID int64, securityID int64)(error){
	sqlStatement := `UPDATE entities SET security = ? WHERE id = ?`

	_, err := db.Exec(sqlStatement, securityID, entityID)

	if err != nil{
		return err
	}

	return nil
}

func (db *MySqlDB)UpdateEntityWithAssocUser(entityID int64, assocUserID int64)(error){
	sqlStatement := `UPDATE entities SET assocUser = ? WHERE id = ?`

	_, err := db.Exec(sqlStatement, assocUserID, entityID)

	if err != nil{
		return err
	}

	return nil
}

func (db *MySqlDB)CompleteEntityIPO(ipoTransactionID, entityID int64)(error){
	sqlStatement := `UPDATE entities SET ipo = ? WHERE id = ?`

	_, err := db.Exec(sqlStatement, ipoTransactionID, entityID)

	if err != nil{
		return err
	}

	return nil
}

func (db *MySqlDB)GetEntitiesWithPrefix(prefix string)([]model.Entity, error){
	sqlStatement := `SELECT * FROM entities WHERE name LIKE CONCAT(?, '%')`

	rows, err := db.Query(sqlStatement, prefix)
	if err != nil{
		return []model.Entity{}, err
	}
	defer rows.Close()

	entities := make([]model.Entity, 0)

	for rows.Next(){
		entity, err := ScanEntity(rows)
		if err != nil{
			return []model.Entity{}, err
		}
		entities = append(entities, *entity)
	}
	return entities, nil
}