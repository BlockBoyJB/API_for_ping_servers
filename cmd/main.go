package main

import (
	"API_for_ping_servers/internal/api"
	"API_for_ping_servers/internal/database"
	"API_for_ping_servers/internal/services"
	"API_for_ping_servers/internal/workerpool"
	"log"
	"net/http"
	"time"
)

var (
	WorkersCount    = 10
	ResponseTimeout = 1 * time.Minute
	RequestInterval = 5 * time.Minute
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	repo := database.New(db)
	if err = repo.CreateTables(); err != nil {
		log.Fatal(err)
	}

	result := make(chan workerpool.JobResult)

	wp := workerpool.NewWP(WorkersCount, ResponseTimeout, result)

	jobs, err := workerpool.GetJobs(repo.PingData)
	if err != nil {
		log.Fatal(err)
	}
	wp.StartWorkers()

	go initJobs(&jobs, wp) // neccessary to pass pointer otherwise the new urls will not be included
	go logResults(result)

	ps := services.PingService{Repo: repo, Wp: wp, Jobs: &jobs}

	ph := api.PingHandler{PingService: ps}

	http.HandleFunc("/api_key", api.GetApiKey)
	http.HandleFunc("/ping/add", ph.AddPing)
	http.HandleFunc("/ping/delete", ph.DeletePing)
	http.HandleFunc("/pings", ph.GetAllPigs)

	if err = http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func initJobs(jobs *[]workerpool.Job, wp *workerpool.Pool) {
	for {
		for _, job := range *jobs {
			wp.AddJob(job)
		}
		time.Sleep(RequestInterval)
	}
}

func logResults(results chan workerpool.JobResult) {
	go func() {
		for result := range results {
			if result.Error != nil {
				log.Printf("[ERROR] | server is not responding to requests | %s", result.Url)
			} else {
				log.Printf("[OK] | Url: %s | Response Duration: %s | Status Code: %d", result.Url, result.ResponseDuration, result.StatusCode)
			}
		}
	}()
}
