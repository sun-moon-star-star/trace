package mysql_test

import (
	"testing"
	"trace/storage"
	"trace/storage/mysql"
)

func TestStatistics(t *testing.T) {
	s := storage.GetStatisticsAPIUseful(&mysql.Mysql{})
	t.Log(s.CountFieldDay("tag", "A-Tag"))
}
