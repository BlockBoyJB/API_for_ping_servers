package workerpool

import (
	"API_for_ping_servers/internal/database"
	"sync"
	"time"
)

type (
	Job struct {
		Url   string
		Email string
	}
	JobResult struct {
		Url              string
		ResponseDuration time.Duration
		StatusCode       int
		Error            error
	}
)

type Pool struct {
	worker  *worker
	wCount  int
	jobs    chan Job
	results chan JobResult
	wg      *sync.WaitGroup
	stopped bool
}

func NewWP(wCount int, timeout time.Duration, result chan JobResult) *Pool {
	return &Pool{
		worker:  createWorker(timeout),
		wCount:  wCount,
		jobs:    make(chan Job),
		results: result,
		wg:      new(sync.WaitGroup),
		stopped: false,
	}
}

func GetJobs(pd *database.PingData) ([]Job, error) {
	data, err := pd.GetAll()
	if err != nil {
		return nil, err
	}
	jobs := make([]Job, 0)
	for _, j := range data {
		jobs = append(jobs, Job{Url: j.Url, Email: j.Email})
	}
	return jobs, nil
}

func (p *Pool) AddJob(j Job) {
	if p.stopped {
		return
	}
	p.jobs <- j
	p.wg.Add(1)
}

func (p *Pool) RemoveJob(j *Job) {
	if p.stopped {
		return
	}

}

func (p *Pool) Stop() {
	p.stopped = true
	close(p.jobs)
	p.wg.Wait()
}

func (p *Pool) initWorker() {
	for job := range p.jobs {
		time.Sleep(time.Second)
		p.results <- p.worker.ping(job)
		p.wg.Done()
	}
}

func (p *Pool) StartWorkers() {
	for i := 0; i < p.wCount; i++ {
		go p.initWorker()
	}
}
