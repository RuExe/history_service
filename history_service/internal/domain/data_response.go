package domain

import (
	"time"
)

type DataItem struct {
	ID         int
	Time       time.Time
	MetricData *MetricData
}

type MetricData struct {
	SellingProductsTotal int
	SellingProductsOnWB  int
	SellingQuantityOnWB  int
	SellingProductsOnMP  int
	SellingQuantityOnMP  int
	ProductsOnWBStores   int
	QuantityOnWBStores   int
	ProductsOnExStores   int
	QuantityOnExStores   int
}
