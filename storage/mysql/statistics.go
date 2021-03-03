package mysql

import "fmt"

func (mysql *Mysql) CountField(tableName, fieldKey string, lowId, highId uint64) (cnt uint64, err error) {
	if err = setDB(); err != nil {
		return
	}

	primaryKey := tableName + "_id"
	prepared := fmt.Sprintf("field = ? and %s between ? and ? ", primaryKey)

	res := db.Table(tableName).Where(prepared, fieldKey, lowId, highId).Count(&cnt)
	err = res.Error

	return
}
