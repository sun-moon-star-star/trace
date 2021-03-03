package storage

import (
	"time"
	"trace/storage/mysql"
)

type StatisticsAPI interface {
	CountField(tableName, fieldName string, lowId, highId uint64) (uint64, error)
	QueryField(tableName, file string, lowId, highId uint64, data interface{}) error
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

	QueryField(tableName, fieldName string, lowId, highId uint64, data interface{}) error
	QueryFieldYear(tableName, fieldName string, data interface{}) error
	QueryFieldMonth(tableName, fieldName string, data interface{}) error
	QueryFieldHalfMonth(tableName, fieldName string, data interface{}) error
	QueryFieldDay(tableName, fieldName string, data interface{}) error
	QueryFieldHour(tableName, fieldName string, data interface{}) error
	QueryFieldMinute(tableName, fieldName string, data interface{}) error
	QueryFieldSecond(tableName, fieldName string, data interface{}) error
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

func (w *WrapperStatisticsAPI) CountField(tableName, fieldName string, lowId, highId uint64) (uint64, error) {
	return w.base.CountField(tableName, fieldName, lowId, highId)
}

func (w *WrapperStatisticsAPI) QueryField(tableName, fieldName string, lowId, highId uint64, data interface{}) error {
	return w.base.QueryField(tableName, fieldName, lowId, highId, data)
}

func (w *WrapperStatisticsAPI) QueryFieldYear(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(12*30*24*60*60*1000<<22), time, data)
}

func (w *WrapperStatisticsAPI) QueryFieldMonth(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(30*24*60*60*1000<<22), time, data)
}

func (w *WrapperStatisticsAPI) QueryFieldHalfMonth(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(15*24*60*60*1000<<22), time, data)
}

func (w *WrapperStatisticsAPI) QueryFieldDay(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(24*60*60*1000<<22), time, data)
}

func (w *WrapperStatisticsAPI) QueryFieldHour(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(60*60*1000<<22), time, data)
}

func (w *WrapperStatisticsAPI) QueryFieldMinute(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(60*1000<<22), time, data)
}

func (w *WrapperStatisticsAPI) QueryFieldSecond(tableName, fieldName string, data interface{}) error {
	time := uint64(time.Now().UnixNano()/1e6) << 22
	return w.base.QueryField(tableName, fieldName, time-(1000<<22), time, data)
}

func GetStatisticsAPIUseful(base StatisticsAPI) StatisticsAPIUseful {
	return &WrapperStatisticsAPI{
		base: base,
	}
}

func DefaultStatisticsAPI() StatisticsAPI {
	return &mysql.Mysql{}
}
