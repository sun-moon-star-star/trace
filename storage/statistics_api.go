package storage

import (
	"time"
	"trace/storage/mysql"
)

type StatisticsAPI interface {
	CountField(tableName, fieldName string, lowId, highId uint64) (uint64, error)
}

type StatisticsAPIUseful interface {
	CountField(tableName, fieldName string, lowId, highId uint64) (uint64, error)
	CountFieldYear(tableName, fieldName string) (uint64, error)      // 365 days
	CountFieldMonth(tableName, fieldName string) (uint64, error)     // 30 days
	CountFieldHalfMonth(tableName, fieldName string) (uint64, error) // 15 days
	CountFieldDay(tableName, fieldName string) (uint64, error)
	CountFieldHour(tableName, fieldName string) (uint64, error)
	CountFieldMinute(tableName, fieldName string) (uint64, error)
	CountFieldSecond(tableName, fieldName string) (uint64, error)
}

type WrapperStatisticsAPI struct {
	base StatisticsAPI
}

func (w *WrapperStatisticsAPI) CountField(tableName, fieldName string, lowId, highId uint64) (uint64, error) {
	return w.base.CountField(tableName, fieldName, lowId, highId)
}

func (w *WrapperStatisticsAPI) CountFieldYear(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(12*30*24*60*60*1000<<22), time)
}

func (w *WrapperStatisticsAPI) CountFieldMonth(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(30*24*60*60*1000<<22), time)
}

func (w *WrapperStatisticsAPI) CountFieldHalfMonth(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(15*24*60*60*1000<<22), time)
}

func (w *WrapperStatisticsAPI) CountFieldDay(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(24*60*60*1000<<22), time)
}

func (w *WrapperStatisticsAPI) CountFieldHour(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(60*60*1000<<22), time)
}

func (w *WrapperStatisticsAPI) CountFieldMinute(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(60*1000<<22), time)
}

func (w *WrapperStatisticsAPI) CountFieldSecond(tableName, fieldName string) (uint64, error) {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.CountField(tableName, fieldName, time-(1000<<22), time)
}

func GetStatisticsAPIUseful(base StatisticsAPI) StatisticsAPIUseful {
	return &WrapperStatisticsAPI{
		base: base,
	}
}

func DefaultStatisticsAPI() StatisticsAPI {
	return &mysql.Mysql{}
}
