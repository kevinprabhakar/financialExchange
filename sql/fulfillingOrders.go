package sql

func (db *MySqlDB) CreateFulfillingOrdersTable()(error){
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS fulfillingOrders(" +
		"id INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"orderPlaced integer," +
		"orderFulfilling integer," +
		"PRIMARY KEY (id) )")
	if err != nil{
		return err
	}
	return nil
}

func (db *MySqlDB) InsertFulfillingOrdersToTable(mainOrder int64, fulfillingOrders []int64)(error){
	query := `INSERT INTO fulfillingOrders ( orderPlaced, orderFulfilling) VALUES ( ?, ?)`

	for _, order := range(fulfillingOrders){
		_, err := db.Exec(query, mainOrder, order)

		if err != nil{
			return err
		}
	}
	return nil
}
