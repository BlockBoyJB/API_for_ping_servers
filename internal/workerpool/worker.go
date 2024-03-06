package workerpool

import (
	"API_for_ping_servers/internal/smtp"
	"net/http"
	"time"
)

type worker struct {
	client *http.Client
}

func createWorker(timeout time.Duration) *worker {
	return &worker{&http.Client{Timeout: timeout}}
}

func (w worker) ping(j Job) JobResult {
	r := JobResult{Url: j.Url}

	start := time.Now()

	resp, err := w.client.Get(j.Url)
	if err != nil {
		r.Error = err
		go smtp.SendPingErrorMail(j.Url, j.Email)
		return r
	}
	r.StatusCode = resp.StatusCode
	r.ResponseDuration = time.Since(start)
	return r
}
