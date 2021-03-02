package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var pool map[string]*gorm.DB

func getConnByDesc(db_desc string) (*gorm.DB, error) {
	if pool == nil {
		pool = make(map[string]*gorm.DB)
	}

	if db, ok := pool[db_desc]; ok {
		return db, nil
	}

	db, err := gorm.Open("mysql", db_desc)
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(time.Duration(3600) * time.Second)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)

	pool[db_desc] = db
	return db, nil
}

func getConn(params map[string]interface{}) (*gorm.DB, error) {
	hostname := params["hostname"]
	port := params["port"]
	username := params["username"]
	password := params["password"]
	network := params["network"]
	database := params["database"]

	db_desc := fmt.Sprintf("%v:%v@%v(%v:%v)/%v", username, password, network, hostname, port, database)
	return getConnByDesc(db_desc)
}

func LoadTables(params map[string]interface{}) ([]string, error) {
	db, err := getConn(params)

	if err != nil {
		return nil, err
	}

	res, err := db.DB().Query("SHOW TABLES")

	if err != nil {
		return nil, err
	}

	defer res.Close()

	var table string
	var tables []string

	for res.Next() {
		res.Scan(&table)
		tables = append(tables, table)
	}

	return tables, nil

}

func SelectTableLimit(params map[string]interface{}) (interface{}, error) {
	db, err := getConn(params)

	if err != nil {
		return nil, err
	}

	table := params["table"].(string)
	num := params["num"]

	data := GetTableArrayType(table)
	if data == nil {
		return nil, err
	}
	res := db.Table(table).Find(data).Limit(num)

	if res.Error != nil {
		return nil, res.Error
	}

	return data, nil
}

func UpdateTableField(params map[string]interface{}) error {
	db, err := getConn(params)

	if err != nil {
		return err
	}

	table := params["table"].(string)
	unique_field := params["unique_field"].(string)
	unique_field_value := params["unique_field_value"]
	field := params["field"]
	value := params["value"]

	data := GetTableType(table)
	if data == nil {
		return errors.New("table not exists in table_reflect")
	}

	res := db.Table(table).Where(unique_field+" = ?", unique_field_value).Update(field, value)

	return res.Error
}

func InsertTable(params map[string]interface{}) (interface{}, error) {
	db, err := getConn(params)

	if err != nil {
		return nil, err
	}

	table := params["table"].(string)
	object_jsonstr := params["object_jsonstr"].(string) // json string

	data := GetTableType(table)
	if data == nil {
		return nil, errors.New("table not exists in table_reflect")
	}

	err = json.Unmarshal([]byte(object_jsonstr), data)

	res := db.Table(table).Create(data)

	if res.Error != nil {
		return nil, res.Error
	}

	return res, nil
}
