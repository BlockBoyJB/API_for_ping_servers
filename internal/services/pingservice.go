package services

import (
	"API_for_ping_servers/internal/database"
	"API_for_ping_servers/internal/workerpool"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"strings"
)

type PingService struct {
	Repo *database.Repository
	Wp   *workerpool.Pool
	Jobs *[]workerpool.Job
}

var (
	data, _      = godotenv.Read(".env")
	SecretApiKey = data["SECRET"]
)

func remove(slice []workerpool.Job, el workerpool.Job) []workerpool.Job {
	for i, job := range slice {
		if job == el {
			slice[i] = slice[len(slice)-1]
			return slice[:len(slice)-1]
		}
	}
	return nil
}

func (p *PingService) CreateNewJob(apiKey string, email string, url string) error {
	if err := p.Repo.PingData.Create(apiKey, email, url); err != nil {
		return err
	}
	job := workerpool.Job{
		Url:   url,
		Email: email,
	}
	*p.Jobs = append(*p.Jobs, job)
	fmt.Println(*p.Jobs)
	p.Wp.AddJob(job)
	return nil
}

func (p *PingService) DeleteJob(apiKey string, url string) error {
	job, err := p.Repo.PingData.Get(apiKey, url)
	if err != nil {
		return err
	}

	if err = p.Repo.PingData.Delete(apiKey, url); err != nil {
		return err
	}
	*p.Jobs = remove(*p.Jobs, workerpool.Job{
		Url:   job.Url,
		Email: job.Email,
	})
	return nil
}

func (p *PingService) GetAllJobs(apiKey string) ([]string, error) {
	rows, err := p.Repo.PingData.GetMany(apiKey)
	if err != nil {
		return nil, err
	}
	r := make([]string, 0)
	for _, row := range rows {
		r = append(r, row.Url)
	}
	return r, nil
}

func CreateApiKey() string {
	salt := hex.EncodeToString([]byte(uuid.NewString()))
	res := sha256.Sum256([]byte(salt + SecretApiKey + "Hello world!"))
	return fmt.Sprintf("%x:%s", res, salt)
}

func VerifyApiKey(apiKey string) bool {
	data := strings.Split(apiKey, ":")
	key, salt := data[0], data[1]
	res := sha256.Sum256([]byte(salt + SecretApiKey + "Hello world!"))
	return key == fmt.Sprintf("%x", res)
}
