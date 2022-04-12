package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"page/internal/config"
	"page/internal/repository"
	"page/internal/service"
)

func main() {
	conf := config.NewConfig()
	log.Print("Server started")
	log.Printf("Server port: %q", conf.Server.Port)

	dataRepository := repository.NewDataRepository(conf)
	dataService := service.NewDataService(dataRepository)
	server := NewServer(conf.Server.Port, dataService)
	server.Start()
}

type Server struct {
	Port        string
	DataService *service.Data
}

func NewServer(port string, dataService *service.Data) *Server {
	return &Server{
		Port:        port,
		DataService: dataService,
	}
}

func (s *Server) Start() {
	s.configureRoutes()

	log.Fatal(http.ListenAndServe(s.Port, nil))
}

func (s *Server) configureRoutes() {
	http.HandleFunc("/data", s.getData)
}

func (s *Server) getData(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	query := r.URL.Query()
	startStr, endStr := query.Get("start"), query.Get("end")
	if startStr == "" || endStr == "" {
		http.Error(w, "query params error", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		log.Print(err)
		http.Error(w, "query params error", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		log.Print(err)
		http.Error(w, "query params error", http.StatusBadRequest)
		return
	}

	data, err := s.DataService.GetDataListByDateRange(start, end)
	if err != nil {
		log.Println(err)
		http.Error(w, "query params error", http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
}
