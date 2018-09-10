package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type MySqlDB struct{
	*sql.DB
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

//Make sure to defer db closing
func OpenDatabase(name string)(*MySqlDB, error){
	db, err := sql.Open("mysql",
						fmt.Sprintf("root:Heb1Pet!@/%s",name))
	if err != nil {
		return nil, err
	}

	return &MySqlDB{db}, nil
}

