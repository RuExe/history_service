package service

import (
	"fmt"
	"time"

	"page/internal/domain"
	"page/internal/repository"
)

type Data struct {
	repository *repository.Data
}

func NewDataService(repository *repository.Data) *Data {
	return &Data{
		repository: repository,
	}
}

func (d *Data) AddData(data *domain.DataItem) (int, error) {
	id, err := d.repository.Add(data)
	if err != nil {
		return -1, fmt.Errorf("DataService AddData error: %q", err)
	}
	return id, nil
}

func (d *Data) GetDataListByDateRange(start, end time.Time) ([]*domain.DataItem, error) {
	list, err := d.repository.GetListByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("DataService GetListByDateRange error: %q", err)
	}
	return list, nil
}
