package main

import (
	"log"
	"time"

	"page/internal/client/metric_service"
	"page/internal/config"
	"page/internal/domain"
	"page/internal/repository"
	"page/internal/service"
)

func main() {
	conf := config.NewConfig()
	log.Print("Parser started")
	log.Printf("Parse URL: %q", conf.Parser.ParseUrl)
	log.Printf("Parse inteval: %q", conf.Parser.Interval)

	client := metricservice.NewClient(conf.Parser.ParseUrl)
	dataRepository := repository.NewDataRepository(conf)
	dataService := service.NewDataService(dataRepository)

	for {
		data, err := client.GetData()
		if err != nil {
			log.Printf("Parser: %q", err)
			continue
		}

		dataItem := &domain.DataItem{
			Time:       time.Now(),
			MetricData: data,
		}

		addedId, err := dataService.AddData(dataItem)
		if err != nil {
			log.Printf("Parser: %q", err)
			continue
		}

		log.Printf("Added data with id-%v :%#v", addedId, dataItem)
		time.Sleep(conf.Parser.Interval)
	}
}
