package config

import (
	"time"
)

type Config struct {
	Server
	Parser
	DataBase
}

func NewConfig() *Config {
	return &Config{
		Server{Port: ":8080"},
		Parser{
			ParseUrl: "http://stats-svc.wbx-ru.svc.k8s.3dcat/prod-counters",
			Interval: 5 * time.Minute,
		},
		DataBase{
			Url:    "user=postgres password=11111 dbname=serviceDB sslmode=disable host=localhost port=5433",
			Driver: "postgres",
		},
	}
}

type Server struct {
	Port string
}

type Parser struct {
	ParseUrl string
	Interval time.Duration
}

type DataBase struct {
	Url    string
	Driver string
}
