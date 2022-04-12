package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"page/internal/config"
	"page/internal/domain"
)

type Data struct {
	config *config.Config
	db     *sqlx.DB
}

func NewDataRepository(config *config.Config) *Data {
	return &Data{
		config: config,
	}
}

func (r *Data) Add(item *domain.DataItem) (int, error) {
	db, err := sqlx.Connect(r.config.DataBase.Driver, r.config.DataBase.Url)
	if err != nil {
		return -1, fmt.Errorf("can't connect to Database %w", err)
	}
	defer func() {
		closeDbError := db.Close()
		if err == nil && closeDbError != nil {
			err = closeDbError
		}
	}()

	tx, err := db.Beginx()
	if err != nil {
		return -1, fmt.Errorf("can't negin transaction %w", err)
	}

	metricDataJson, err := json.Marshal(item.MetricData)
	if err != nil {
		return -1, fmt.Errorf("can't marshal data")
	}

	var id int
	query := `INSERT INTO data (json_data, timestamp) VALUES ($1, $2) RETURNING id`
	if err = tx.QueryRowx(query,
		metricDataJson,
		item.Time.UTC(),
	).Scan(&id); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		return -1, fmt.Errorf("can't commit transaction %w", err)
	}
	return id, nil
}

func (r *Data) GetListByDateRange(start, end time.Time) ([]*domain.DataItem, error) {
	db, err := sqlx.Connect(r.config.DataBase.Driver, r.config.DataBase.Url)
	if err != nil {
		return nil, fmt.Errorf("can't connect to Database %w", err)
	}
	defer func() {
		closeDbError := db.Close()
		if err == nil && closeDbError != nil {
			err = closeDbError
		}
	}()

	temps := make([]*temp, 0)

	query := fmt.Sprintf(`SELECT json_data, timestamp FROM data WHERE timestamp >= '%v' AND timestamp <= '%v' ORDER BY timestamp`, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err := db.Select(&temps, query); err != nil {
		return nil, err
	}

	items := make([]*domain.DataItem, len(temps))
	for index, item := range temps {
		data := new(domain.MetricData)
		json.Unmarshal(item.JsonData, data)

		items[index] = &domain.DataItem{
			Time:       item.TimeStamp,
			MetricData: data,
		}
	}
	return items, nil
}

type temp struct {
	TimeStamp time.Time `db:"timestamp"`
	JsonData  []byte    `db:"json_data"`
}
